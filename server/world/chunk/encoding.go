package chunk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type (
	// Encoding is an encoding type used for Chunk encoding. Implementations of this interface are DiskEncoding and
	// NetworkEncoding, which can be used to encode a Chunk to an intermediate disk or network representation respectively.
	Encoding interface {
		encodePalette(buf *bytes.Buffer, p *Palette, e paletteEncoding)
		decodePalette(buf *bytes.Buffer, blockSize paletteSize, e paletteEncoding) (*Palette, error)
		network() byte
	}
	// paletteEncoding is an encoding type used for Chunk encoding. It is used to encode different types of palettes
	// (for example, blocks or biomes) differently.
	paletteEncoding interface {
		encode(buf *bytes.Buffer, v uint32)
		decode(buf *bytes.Buffer, e Encoding) (uint32, error)
		withCustomBlocks() bool
	}
)

var (
	// DiskEncoding is the Encoding for writing a Chunk to disk. It writes block palettes using NBT and does not use
	// varints.
	DiskEncoding diskEncoding
	// NetworkEncoding is the Encoding used for sending a Chunk over network. It does not use NBT and writes varints.
	NetworkEncoding networkEncoding

	// UnknownRID is the runtime ID (IN VANILLA) of the unknown block.
	UnknownRID uint32
)

// BiomePaletteEncoding implements the encoding of biome palettes to disk.
type BiomePaletteEncoding struct{}

func (BiomePaletteEncoding) encode(buf *bytes.Buffer, v uint32) {
	_ = binary.Write(buf, binary.LittleEndian, v)
}

func (BiomePaletteEncoding) decode(buf *bytes.Buffer, e Encoding) (uint32, error) {
	var v uint32
	return v, binary.Read(buf, binary.LittleEndian, &v)
}

func (e BiomePaletteEncoding) withCustomBlocks() bool {
	return false
}

// BlockPaletteEncoding implements the encoding of block palettes to disk.
type BlockPaletteEncoding struct {
	HasCustom bool
}

func (e BlockPaletteEncoding) withCustomBlocks() bool {
	return e.HasCustom
}

func (BlockPaletteEncoding) encode(buf *bytes.Buffer, v uint32) {
	// Get the block state registered with the runtime IDs we have in the palette of the block storage
	// as we need the name and data value to store.
	name, props, _ := RuntimeIDToState(v)
	_ = nbt.NewEncoderWithEncoding(buf, nbt.LittleEndian).Encode(blockEntry{Name: name, State: props, Version: CurrentBlockVersion})
}

func (BlockPaletteEncoding) decode(buf *bytes.Buffer, e Encoding) (uint32, error) {
	var m map[string]any

	var nbt_e nbt.Encoding = nbt.LittleEndian
	if e.network() == 1 {
		nbt_e = nbt.NetworkLittleEndian
	}

	if err := nbt.NewDecoderWithEncoding(buf, nbt_e).Decode(&m); err != nil {
		return 0, fmt.Errorf("error decoding block palette entry: %w", err)
	}

	// Decode the name and version of the block entry.
	name, _ := m["name"].(string)
	version, _ := m["version"].(int32)

	// Now check for a state field.
	stateI, ok := m["states"]
	if !ok {
		// If it doesn't exist, this is likely a pre-1.13 block state, so decode the meta value instead.
		meta, _ := m["val"].(int16)

		// Upgrade the pre-1.13 state into a post-1.13 state.
		state, ok := upgradeLegacyEntry(name, meta)
		if !ok {
			return 0, fmt.Errorf("cannot find mapping for legacy block entry: %v, %v", name, meta)
		}

		// Update the state.
		stateI = state.State
	}
	state, ok := stateI.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("invalid state in block entry")
	}

	if !strings.Contains(name, ":") { // why is this needed?
		name = "minecraft:" + name
	}

	// If the entry is an alias, then we need to resolve it.
	entry := blockEntry{Name: name, State: state, Version: version}
	if updatedEntry, ok := upgradeAliasEntry(entry); ok {
		entry = updatedEntry
	}

	v, ok := StateToRuntimeID(entry.Name, entry.State)
	if !ok {
		return 0, fmt.Errorf("cannot get runtime ID of block state %v{%+v}", name, state)
	}
	return v, nil
}

// diskEncoding implements the Chunk encoding for writing to disk.
type diskEncoding struct{}

func (diskEncoding) network() byte { return 0 }
func (diskEncoding) encodePalette(buf *bytes.Buffer, p *Palette, e paletteEncoding) {
	if p.size != 0 {
		_ = binary.Write(buf, binary.LittleEndian, uint32(p.Len()))
	}
	for _, v := range p.values {
		e.encode(buf, v)
	}
}

func (diskEncoding) decodePalette(buf *bytes.Buffer, blockSize paletteSize, e paletteEncoding) (*Palette, error) {
	paletteCount := uint32(1)
	if blockSize != 0 {
		if err := binary.Read(buf, binary.LittleEndian, &paletteCount); err != nil {
			return nil, fmt.Errorf("error reading palette entry count: %w", err)
		}
	}

	var err error
	palette := newPalette(blockSize, make([]uint32, paletteCount))
	for i := uint32(0); i < paletteCount; i++ {
		palette.values[i], err = e.decode(buf, DiskEncoding)
		if err != nil {
			return nil, err
		}
	}
	if paletteCount == 0 {
		return palette, fmt.Errorf("invalid palette entry count: found 0, but palette with %v bits per block must have at least 1 value", blockSize)
	}
	return palette, nil
}

// networkEncoding implements the Chunk encoding for sending over network.
type networkEncoding struct{}

func (networkEncoding) network() byte { return 1 }
func (networkEncoding) encodePalette(buf *bytes.Buffer, p *Palette, _ paletteEncoding) {
	if p.size != 0 {
		_ = protocol.WriteVarint32(buf, int32(p.Len()))
	}
	for _, val := range p.values {
		_ = protocol.WriteVarint32(buf, int32(val))
	}
}

func (networkEncoding) decodePalette(buf *bytes.Buffer, blockSize paletteSize, e paletteEncoding) (*Palette, error) {
	var paletteCount int32 = 1
	if blockSize != 0 {
		if err := protocol.Varint32(buf, &paletteCount); err != nil {
			return nil, fmt.Errorf("error reading palette entry count: %w", err)
		}
		if paletteCount <= 0 {
			return nil, fmt.Errorf("invalid palette entry count %v", paletteCount)
		}
	}

	magic := buf.Bytes()[0:2]
	is_legacy := bytes.Equal(magic, []byte{0x0A, 0x00})

	if is_legacy { // is using nbt data
		var err error
		palette := newPalette(blockSize, make([]uint32, paletteCount))
		for i := int32(0); i < paletteCount; i++ {
			palette.values[i], err = e.decode(buf, NetworkEncoding)
			if err != nil {
				return nil, err
			}
		}
		return palette, nil
	} else { // regular block rid list
		blocks, temp := make([]uint32, paletteCount), int32(0)

		for i := int32(0); i < paletteCount; i++ {
			if err := protocol.Varint32(buf, &temp); err != nil {
				return nil, fmt.Errorf("error decoding palette entry: %w", err)
			}
			rid := uint32(temp)
			if e.withCustomBlocks() {
				if rid > UnknownRID {
					rid -= 1
				}
			}
			blocks[i] = rid
		}

		return &Palette{values: blocks, size: blockSize}, nil
	}
}
