package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	out := flag.String("o", "", "output file for hash constants and methods")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalln("Must pass one package to produce block hashes for.")
	}
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedFiles,
	}
	pkgs, err := packages.Load(cfg, flag.Args()[0])
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.OpenFile(*out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	for _, pkg := range pkgs {
		procPackage(pkg, f)
	}
	_ = f.Close()
}

func procPackage(pkg *packages.Package, w io.Writer) {
	b := &hashBuilder{
		pkg:         pkg,
		fields:      make(map[string][]*ast.Field),
		aliases:     make(map[string]string),
		handled:     map[string]struct{}{},
		funcs:       map[string]*ast.FuncDecl{},
		blockFields: map[string][]*ast.Field{},
	}
	b.readStructFields(pkg)
	b.readFuncs(pkg)
	b.resolveBlocks()
	b.sortNames()

	b.writePackage(w)
	b.writeConstants(w)
	b.writeNextHash(w)
	b.writeMethods(w)
}

var (
	packageFormat = "// Code generated by cmd/blockhash; DO NOT EDIT.\n\npackage %v\n\n"
	methodFormat  = "\nfunc (%v%v) Hash() (uint64, uint64) {\n\treturn %v, %v\n}\n"
	constFormat   = "\thash%v"
)

type hashBuilder struct {
	pkg         *packages.Package
	fields      map[string][]*ast.Field
	funcs       map[string]*ast.FuncDecl
	aliases     map[string]string
	handled     map[string]struct{}
	blockFields map[string][]*ast.Field
	names       []string
}

// sortNames sorts the names of the blockFields map and stores them in a slice.
func (b *hashBuilder) sortNames() {
	b.names = make([]string, 0, len(b.blockFields))
	for name := range b.blockFields {
		b.names = append(b.names, name)
	}
	sort.Slice(b.names, func(i, j int) bool {
		return b.names[i] < b.names[j]
	})
}

// writePackage writes the package at the top of the file.
func (b *hashBuilder) writePackage(w io.Writer) {
	if _, err := fmt.Fprintf(w, packageFormat, b.pkg.Name); err != nil {
		log.Fatalln(err)
	}
	if _, err := fmt.Fprintf(w, "import \"github.com/df-mc/dragonfly/server/world\"\n\n"); err != nil {
		log.Fatalln(err)
	}
}

// writeConstants writes hash constants for every block to a file.
func (b *hashBuilder) writeConstants(w io.Writer) {
	if _, err := fmt.Fprintln(w, "const ("); err != nil {
		log.Fatalln(err)
	}

	for i, name := range b.names {
		c := constFormat
		if i == 0 {
			c += " = iota"
		}

		if _, err := fmt.Fprintf(w, c+"\n", name); err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := fmt.Fprintln(w, "\thashCustomBlockBase\n)"); err != nil {
		log.Fatalln(err)
	}
}

func (b *hashBuilder) writeNextHash(w io.Writer) {
	if _, err := fmt.Fprintln(w, "\n// customBlockBase represents the base hash for all custom blocks."); err != nil {
		log.Fatalln(err)
	}
	if _, err := fmt.Fprintln(w, "var customBlockBase = uint64(hashCustomBlockBase - 1)"); err != nil {
		log.Fatalln(err)
	}
	if _, err := fmt.Fprintln(w, "\n// NextHash returns the next free hash for custom blocks."); err != nil {
		log.Fatalln(err)
	}
	if _, err := fmt.Fprintln(w, "func NextHash() uint64 {\n\tcustomBlockBase++\n\treturn customBlockBase\n}"); err != nil {
		log.Fatalln(err)
	}
}

func (b *hashBuilder) writeMethods(w io.Writer) {
	for _, name := range b.names {
		fields := b.blockFields[name]

		var h string
		var bitSize int

		fun := b.funcs[name]
		var recvName string
		for _, n := range fun.Recv.List[0].Names {
			recvName = n.Name
		}
		pos := b.pkg.Fset.Position(fun.Body.Pos())
		f, err := os.Open(pos.Filename)
		if err != nil {
			log.Fatalln(err)
		}
		body := make([]byte, fun.Body.End()-fun.Body.Pos())

		if _, err := f.ReadAt(body, int64(pos.Offset)); err != nil {
			log.Fatalln(err)
		}
		_ = f.Close()

		for _, field := range fields {
			for _, fieldName := range field.Names {
				if !bytes.Contains(body, []byte(fieldName.Name)) {
					// Field was not used in the EncodeBlock method, so we can assume it's not a property and thus
					// should not be in the Hash method.
					continue
				}
				if !fieldName.IsExported() {
					continue
				}
				directives := make(map[string]string)
				if field.Doc != nil {
					for _, d := range field.Doc.List {
						const k = "//blockhash:"
						if index := strings.Index(d.Text, k); index != -1 {
							dir := strings.Split(d.Text[index+len(k):], " ")
							directives[dir[0]] = strings.Join(dir[1:], " ")
						}
					}
				}
				str, v := b.ftype(name, recvName+"."+fieldName.Name, field.Type, directives)
				if v == 0 {
					// Assume this field is not used in the hash.
					continue
				}

				if bitSize > 64 {
					log.Println("Hash size of block properties of", name, "exceeds 64 bits. Please look at this manually.")
				} else {
					if h == "" {
						h += str
					} else {
						h += " | " + str
					}
					if bitSize > 0 {
						h += "<<" + strconv.Itoa(bitSize)
					}
				}
				bitSize += v
			}
		}
		if bitSize == 0 {
			// No need to have a receiver name if we don't use any of the fields of the block.
			recvName = ""
		}

		if recvName != "" {
			recvName += " "
		}
		if h == "" {
			h = "0"
		}

		if _, err := fmt.Fprintf(w, methodFormat, recvName, name, "hash"+name, h); err != nil {
			log.Fatalln(err)
		}
	}
}

func (b *hashBuilder) ftype(structName, s string, expr ast.Expr, directives map[string]string) (string, int) {
	var name string
	switch t := expr.(type) {
	case *ast.BasicLit:
		name = t.Value
	case *ast.Ident:
		name = t.Name
	case *ast.SelectorExpr:
		name = t.Sel.Name
	default:
		log.Fatalf("unknown field type %#v\n", expr)
		return "", 0
	}
	switch name {
	case "bool":
		return "uint64(boolByte(" + s + "))", 1
	case "int":
		return "uint64(" + s + ")", 8
	case "Block":
		return "world.BlockHash(" + s + ")", 32
	case "Attachment":
		if _, ok := directives["facing_only"]; ok {
			log.Println("Found directive: 'facing_only'")
			return "uint64(" + s + ".FaceUint8())", 3
		}
		return "uint64(" + s + ".Uint8())", 5
	case "GrindstoneAttachment":
		return "uint64(" + s + ".Uint8())", 2
	case "WoodType", "FlowerType", "DoubleFlowerType", "Colour", "MushroomType", "SeagrassType":
		// Assuming these were all based on metadata, it should be safe to assume a bit size of 4 for this.
		return "uint64(" + s + ".Uint8())", 4
	case "CoralType", "SkullType":
		return "uint64(" + s + ".Uint8())", 3
	case "AnvilType", "SandstoneType", "PrismarineType", "StoneBricksType", "NetherBricksType", "FroglightType",
		"WallConnectionType", "BlackstoneType", "DeepslateType", "TallGrassType", "CopperType", "OxidationType":
		return "uint64(" + s + ".Uint8())", 2
	case "OreType", "FireType", "DoubleTallGrassType":
		return "uint64(" + s + ".Uint8())", 1
	case "Direction", "Axis":
		return "uint64(" + s + ")", 2
	case "Face":
		return "uint64(" + s + ")", 3
	default:
		log.Println("Found unhandled field type", "'"+name+"'", "in block", structName+".", "Assuming this field is not included in block states. Please make sure this is correct or add the type to cmd/blockhash.")
	}
	return "", 0
}

func (b *hashBuilder) resolveBlocks() {
	for bl, fields := range b.fields {
		if _, ok := b.funcs[bl]; ok {
			b.blockFields[bl] = fields
		}
	}
}

func (b *hashBuilder) readFuncs(pkg *packages.Package) {
	for _, f := range pkg.Syntax {
		ast.Inspect(f, b.readFuncDecls)
	}
}

func (b *hashBuilder) readFuncDecls(node ast.Node) bool {
	if fun, ok := node.(*ast.FuncDecl); ok {
		// If the function is called 'EncodeBlock' and the receiver is not nil, meaning the function is a method, this
		// is an implementation of the world.Block interface.
		if fun.Name.Name == "EncodeBlock" && fun.Recv != nil {
			b.funcs[fun.Recv.List[0].Type.(*ast.Ident).Name] = fun
		}
	}
	return true
}

func (b *hashBuilder) readStructFields(pkg *packages.Package) {
	for _, f := range pkg.Syntax {
		ast.Inspect(f, b.readStructs)
	}
	b.resolveEmbedded()
	b.resolveAliases()
}

func (b *hashBuilder) resolveAliases() {
	for name, alias := range b.aliases {
		b.fields[name] = b.findFields(alias)
	}
}

func (b *hashBuilder) findFields(structName string) []*ast.Field {
	for {
		if fields, ok := b.fields[structName]; ok {
			// Alias found in the fields map, so it referred to a struct directly.
			return fields
		}
		if nested, ok := b.aliases[structName]; ok {
			// The alias itself was an alias, so continue with the next.
			structName = nested
			continue
		}
		// Neither an alias nor a struct: Break as this isn't going to go anywhere.
		return nil
	}
}

func (b *hashBuilder) resolveEmbedded() {
	for name, fields := range b.fields {
		if _, ok := b.handled[name]; ok {
			// Don't handle if a previous run already handled this struct.
			continue
		}
		newFields := make([]*ast.Field, 0, len(fields))
		for _, f := range fields {
			if len(f.Names) == 0 {
				// We're dealing with an embedded struct here. They're of the type ast.Ident.
				if ident, ok := f.Type.(*ast.Ident); ok {
					for _, af := range b.findFields(ident.Name) {
						if len(af.Names) == 0 {
							// The struct this referred is embedding a struct itself which hasn't yet been processed,
							// so we need to rerun and hope that struct is handled next. This isn't a very elegant way,
							// and could lead to a lot of runs, but in general it's fast enough and does the job.
							b.resolveEmbedded()
							return
						}
					}
					newFields = append(newFields, b.findFields(ident.Name)...)
				}
			} else {
				newFields = append(newFields, f)
			}
		}
		// Make sure a next run doesn't end up handling this struct again.
		b.handled[name] = struct{}{}
		b.fields[name] = newFields
	}
}

func (b *hashBuilder) readStructs(node ast.Node) bool {
	if s, ok := node.(*ast.TypeSpec); ok {
		switch t := s.Type.(type) {
		case *ast.StructType:
			b.fields[s.Name.Name] = t.Fields.List
		case *ast.Ident:
			// This is a type created something like 'type Andesite polishable': A type alias. We need to handle
			// these later, first parse all struct types.
			b.aliases[s.Name.Name] = t.Name
		}
	}
	return true
}
