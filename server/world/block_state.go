package world

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"unsafe"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

var (
	//go:embed block_states.nbt
	blockStateData []byte
)

func init() {
	dec := nbt.NewDecoder(bytes.NewBuffer(blockStateData))

	// Register all block states present in the block_states.nbt file. These are all possible options registered
	// blocks may encode to.
	for {
		var s BlockState
		if err := dec.Decode(&s); err != nil {
			break
		}
		DefaultBlockRegistry.RegisterBlockState(s)
	}
}

var bufNetworkhash []byte = make([]byte, 0xff)

func networkBlockHash(name string, properties map[string]any) uint32 {
	if name == "minecraft:unknown" {
		return 0xfffffffe // -2
	}

	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	bufNetworkhash = bufNetworkhash[:0]
	var data = bufNetworkhash
	writeString := func(str string) {
		data = binary.LittleEndian.AppendUint16(data, uint16(len(str)))
		data = append(data, []byte(str)...)
	}

	data = append(data, 10) // compound
	data = append(data, 0)
	data = append(data, 0)

	data = append(data, 8) // string
	writeString("name")
	writeString(name)

	data = append(data, 10) // compound
	writeString("states")
	for _, k := range keys {
		v := properties[k]
		switch v := v.(type) {
		case string:
			data = append(data, 8) // string
			writeString(k)
			writeString(v)

		case uint8:
			data = append(data, 1) // tagByte
			writeString(k)
			data = append(data, byte(v))
		case int8:
			data = append(data, 1) // tagByte
			writeString(k)
			data = append(data, byte(v))
		case bool:
			b := 0
			if v {
				b = 1
			}
			data = append(data, 1) // tagByte
			writeString(k)
			data = append(data, byte(b))

		case uint16:
			data = append(data, 2) // tagInt16
			writeString(k)
			data = binary.LittleEndian.AppendUint16(data, uint16(v))
		case int16:
			data = append(data, 2) // tagInt16
			writeString(k)
			data = binary.LittleEndian.AppendUint16(data, uint16(v))

		case uint32:
			data = append(data, 3) // tagInt32
			writeString(k)
			data = binary.LittleEndian.AppendUint32(data, uint32(v))
		case int32:
			data = append(data, 3) // tagInt32
			writeString(k)
			data = binary.LittleEndian.AppendUint32(data, uint32(v))
		default:
			panic("unhandled nbt type")
		}
	}
	data = append(data, 0) // end
	data = append(data, 0) // end

	h := fnv.New32a()
	h.Write(data)
	return h.Sum32()
}

func splitNamespace(identifier string) (ns, name string) {
	ns_name := strings.Split(identifier, ":")
	return ns_name[0], ns_name[len(ns_name)-1]
}

// BlockState holds a combination of a name and properties, together with a version.
type BlockState struct {
	Name       string         `nbt:"name"`
	Properties map[string]any `nbt:"states"`
	Version    int32          `nbt:"version"`
}

// stateHash is a struct that may be used as a map key for block states. It contains the name of the block state
// and an encoded version of the properties.
type stateHash struct {
	name, properties string
}

// hashProperties produces a hash for the block properties held by the blockState.
var hashPropertiesBuilder strings.Builder

func hashProperties(properties map[string]any) string {
	if properties == nil {
		return ""
	}
	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	hashPropertiesBuilder.Reset()
	var b = hashPropertiesBuilder
	for _, k := range keys {
		switch v := properties[k].(type) {
		case bool:
			if v {
				b.WriteByte(1)
			} else {
				b.WriteByte(0)
			}
		case uint8:
			b.WriteByte(v)
		case int32:
			a := *(*[4]byte)(unsafe.Pointer(&v))
			b.Write(a[:])
		case string:
			b.WriteString(v)
		default:
			// If block encoding is broken, we want to find out as soon as possible. This saves a lot of time
			// debugging in-game.
			panic(fmt.Sprintf("invalid block property type %T for property %v", v, k))
		}
	}

	return b.String()
}
