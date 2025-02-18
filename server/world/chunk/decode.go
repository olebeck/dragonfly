package chunk

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

// NetworkDecode decodes the network serialised data passed into a Chunk if successful. If not, the chunk
// returned is nil and the error non-nil.
// The sub chunk count passed must be that found in the LevelChunk packet.
// noinspection GoUnusedExportedFunction
func NetworkDecode(br BlockRegistry, data []byte, count int, oldBiomes bool, hashedRids bool, r cube.Range) (*Chunk, []map[string]any, error) {
	var (
		c   = New(br, r)
		buf = bytes.NewBuffer(data)
		err error
	)

	for i := 0; i < count; i++ {
		index := uint8(i)
		c.sub[index], err = decodeSubChunk(buf, c, &index, NetworkEncoding, hashedRids)
		if err != nil {
			return nil, nil, err
		}
	}

	err = DecodeNetworkBiomes(c, buf, oldBiomes)
	if err != nil {
		return nil, nil, err
	}

	borderBlocks, _ := buf.ReadByte()
	buf.Next(int(borderBlocks))

	blockNBTs, err := DecodeBlockNBTs(buf)
	if err != nil {
		return nil, nil, err
	}
	return c, blockNBTs, nil
}

func DecodeNetworkBiomes(c *Chunk, buf *bytes.Buffer, oldBiomes bool) error {
	if oldBiomes {
		// Read the old biomes.
		biomes := make([]byte, 256)
		if _, err := buf.Read(biomes[:]); err != nil {
			return fmt.Errorf("error reading biomes: %w", err)
		}
		var values []uint32
		for _, v := range biomes {
			if !slices.Contains(values, uint32(v)) {
				values = append(values, uint32(v))
			}
		}

		size := paletteSizeFor(len(values))
		biome := newPalettedStorage(make([]uint32, size.uint32s()), newPalette(size, values))

		// Make our 2D biomes 3D.
		for x := 0; x < 16; x++ {
			for z := 0; z < 16; z++ {
				id := biomes[(x&15)|(z&15)<<4]
				for y := 0; y < 16; y++ {
					biome.Set(uint8(x), uint8(y), uint8(z), uint32(id))
				}
			}
		}

		for i := range c.biomes {
			c.biomes[i] = biome
		}
	} else {
		var last *PalettedStorage
		for i := 0; i < len(c.sub); i++ {
			b, err := decodePalettedStorage(buf, NetworkEncoding, BiomePaletteEncoding)
			if err != nil {
				return err
			}
			if b == nil {
				// b == nil means this paletted storage had the flag pointing to the previous one. It basically means we should
				// inherit whatever palette we decoded last.
				if i == 0 {
					// This should never happen and there is no way to handle this.
					return fmt.Errorf("first biome storage pointed to previous one")
				}
				b = last
			} else {
				last = b
			}
			c.biomes[i] = b
		}
	}
	return nil
}

func DecodeBlockNBTs(buf *bytes.Buffer) ([]map[string]any, error) {
	var blockNBTs []map[string]any
	if buf.Len() > 0 {
		dec := nbt.NewDecoderWithEncoding(buf, nbt.NetworkLittleEndian)
		dec.AllowZero = true
		for buf.Len() > 0 {
			blockNBT := make(map[string]any, 0)
			err := dec.Decode(&blockNBT)
			if err != nil {
				return nil, err
			}
			if len(blockNBT) > 0 {
				blockNBTs = append(blockNBTs, blockNBT)
			}
		}
	}
	return blockNBTs, nil
}

// DiskDecode decodes the data from a SerialisedData object into a chunk and returns it. If the data was
// invalid, an error is returned.
func DiskDecode(br BlockRegistry, data SerialisedData, r cube.Range) (*Chunk, error) {
	c := New(br, r)

	err := decodeBiomes(bytes.NewBuffer(data.Biomes), c, DiskEncoding)
	if err != nil {
		return nil, err
	}
	for i, sub := range data.SubChunks {
		if len(sub) == 0 {
			// No data for this sub chunk.
			continue
		}
		index := uint8(i)
		if c.sub[index], err = decodeSubChunk(bytes.NewBuffer(sub), c, &index, DiskEncoding, false); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// decodeSubChunk decodes a SubChunk from a bytes.Buffer. The Encoding passed defines how the block storages of the
// SubChunk are decoded.
func decodeSubChunk(buf *bytes.Buffer, c *Chunk, index *byte, e Encoding, hashedRids bool) (*SubChunk, error) {
	ver, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading version: %w", err)
	}
	sub := NewSubChunk(c.BlockRegistry)
	switch ver {
	default:
		return nil, fmt.Errorf("unknown sub chunk version %v: can't decode", ver)
	case 1:
		// Version 1 only has one layer for each sub chunk, but uses the format with palettes.
		storage, err := decodePalettedStorage(buf, e, BlockPaletteEncoding{Blocks: c.BlockRegistry})
		if err != nil {
			return nil, err
		}
		sub.storages = append(sub.storages, storage)
	case 8, 9:
		// Version 8 allows up to 256 layers for one sub chunk.
		storageCount, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("error reading storage count: %w", err)
		}
		if ver == 9 {
			uIndex, err := buf.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("error reading sub-chunk index: %w", err)
			}
			// The index as written here isn't the actual index of the sub-chunk within the chunk. Rather, it is the Y
			// value of the sub-chunk. This means that we need to translate it to an index.
			*index = uint8(int8(uIndex) - int8(c.r[0]>>4))
		}
		sub.storages = make([]*PalettedStorage, storageCount)

		for i := byte(0); i < storageCount; i++ {
			storage, err := decodePalettedStorage(buf, e, BlockPaletteEncoding{Blocks: c.BlockRegistry})
			if err != nil {
				return nil, err
			}

			if hashedRids {
				for i2, v := range storage.palette.values {
					var ok bool
					storage.palette.values[i2], ok = c.BlockRegistry.HashToRuntimeID(v)
					if !ok {
						/* for debug
						for x := byte(0); x < 16; x++ {
							for y := byte(0); y < 16; y++ {
								for z := byte(0); z < 16; z++ {
									idx := storage.paletteIndex(x, y, z)
									if idx == uint16(i2) {
										println()
									}
								}
							}
						}
						*/
						fmt.Println("rid hash not found, data sorting wrong.")
					}
				}
			}

			sub.storages[i] = storage

		}
	}
	return sub, nil
}

func DecodeSubChunk(buf *bytes.Buffer, br BlockRegistry, r cube.Range, index *byte, e Encoding, hashedRids bool) (*SubChunk, error) {
	return decodeSubChunk(buf, &Chunk{BlockRegistry: br, r: r}, index, e, hashedRids)
}

// decodeBiomes reads the paletted storages holding biomes from buf and stores it into the Chunk passed.
func decodeBiomes(buf *bytes.Buffer, c *Chunk, e Encoding) error {
	var last *PalettedStorage
	if buf.Len() != 0 {
		for i := 0; i < len(c.sub); i++ {
			if i == 16 && buf.Len() == 0 { // fix for 255 height worlds
				copy(c.biomes[4:], c.biomes[:16])
				for j := 0; j < 4; j++ {
					c.biomes[j] = emptyStorage(0)
				}
				break
			}
			b, err := decodePalettedStorage(buf, e, BiomePaletteEncoding)
			if err != nil {
				return err
			}
			// b == nil means this paletted storage had the flag pointing to the previous one. It basically means we should
			// inherit whatever palette we decoded last.
			if i == 0 && b == nil {
				// This should never happen and there is no way to handle this.
				return fmt.Errorf("first biome storage pointed to previous one")
			}
			if b == nil {
				// This means this paletted storage had the flag pointing to the previous one. It basically means we should
				// inherit whatever palette we decoded last.
				b = last
			} else {
				last = b
			}
			c.biomes[i] = b
		}
	}
	return nil
}

// decodePalettedStorage decodes a PalettedStorage from a bytes.Buffer. The Encoding passed is used to read either a
// network or disk block storage.
func decodePalettedStorage(buf *bytes.Buffer, e Encoding, pe paletteEncoding) (*PalettedStorage, error) {
	blockSize, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading block size: %w", err)
	}
	_, isNetwork := e.(networkEncoding)
	_, isBlocks := pe.(BlockPaletteEncoding)
	if isNetwork && isBlocks && blockSize&1 != 1 {
		e = NetworkPersistentEncoding
	}

	blockSize >>= 1
	if blockSize == 0x7f {
		return nil, nil
	}

	size := paletteSize(blockSize)
	if size > 32 {
		return nil, fmt.Errorf("cannot read paletted storage (size=%v) %T: size too large", blockSize, pe)
	}
	uint32Count := size.uint32s()

	uint32s := make([]uint32, uint32Count)
	byteCount := uint32Count * 4

	data := buf.Next(byteCount)
	if len(data) != byteCount {
		return nil, fmt.Errorf("cannot read paletted storage (size=%v) %T: not enough block data present: expected %v bytes, got %v", blockSize, pe, byteCount, len(data))
	}
	for i := 0; i < uint32Count; i++ {
		// Explicitly don't use the binary package to greatly improve performance of reading the uint32s.
		uint32s[i] = uint32(data[i*4]) | uint32(data[i*4+1])<<8 | uint32(data[i*4+2])<<16 | uint32(data[i*4+3])<<24
	}
	p, err := e.decodePalette(buf, paletteSize(blockSize), pe)
	return newPalettedStorage(uint32s, p), err
}
