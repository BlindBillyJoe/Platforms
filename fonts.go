package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	statusNone  = 0
	statusLaded = 1
	statusFreed = 2
)

//Table types
const (
	Cmap  = 0x70616d63
	Cvt   = 0x20747663
	Fpgm  = 0x6d677066
	Glyph = 0x66796c67
	Head  = 0x64616568
	Hhea  = 0x61656868
	Hmtx  = 0x78746d68
	Loca  = 0x61636f6c
	Maxp  = 0x7078616d
	Post  = 0x74736f70
)

//Glyph type
const (
	Compound = 0
	Simple   = 1
)

type ttfFont struct {
	fd *os.File

	scalerType    uint32
	numTables     uint16
	searchRange   uint16
	entrySelector uint16
	rangeShift    uint16

	tables []ttfTable

	point       int16
	dpi         uint16
	ppem        uint16
	upem        uint16
	rasterFlags uint32
}

type ttfTableI interface {
	tableType() int
}

type ttfGlyphI interface {
	glyphType() int
}

type ttfTable struct {
	tag      uint32
	checkSum uint32
	offset   uint32
	length   uint32
	status   uint8
	data     ttfTableI
}

type postTable struct {
	format             uint32
	italicAngle        uint32
	underlinePosition  int16
	underlineThickness int16
	isFixedPitch       uint32
	minMemType42       uint32
	maxMemType42       uint32
	minMemType1        uint32
	maxMemType1        uint32

	numGlyphs  uint16
	glyphNames [][]uint8
}

func (t *postTable) tableType() int {
	return Post
}

type maxpTable struct {
	version              uint32
	numGlyphs            uint16
	maxPoints            uint16
	maxContours          uint16
	maxComponentPoints   uint16
	maxComponentContours uint16
	maxZones             uint16
	maxTwilightPoints    uint16
	maxStorage           uint16
	maxFunctionDefs      uint16
	maxInstructionDefs   uint16
	maxStackElements     uint16
	maxSizeOfInstruction uint16
	maxComponentElements uint16
	maxComponentDepth    uint16
}

func (t *maxpTable) tableType() int {
	return Maxp
}

type locaTable struct {
	offsets    []uint32
	numOffsets uint16
}

func (t *locaTable) tableType() int {
	return Loca
}

type hmtxTable struct {
	advanceWidth                 *uint16
	leftSideBearing              *int16
	nonHorizontalLeftSideBearing *int16

	numHMetrics             uint16
	numNonHorizontalMetrics uint16
}

func (t *hmtxTable) tableType() int {
	return Hmtx
}

type hheaTable struct {
	version             uint32
	ascent              int16
	descent             int16
	lineGap             int16
	advanceWidthMax     uint16
	minLeftSideBearing  int16
	minRightSideBearing int16
	xMaxExtent          int16
	caretSlopeRise      int16
	caretSlopeRun       int16
	caretOffset         int16
	reserved1           int16
	reserved2           int16
	reserved3           int16
	reserved4           int16
	metricDataFormat    int16
	numOfLongHorMetrics uint16
}

func (t *hheaTable) tableType() int {
	return Hhea
}

type headTable struct {
	version            uint32
	fontRevision       uint32
	checkSumAdjustment uint32
	magicNumber        uint32
	flags              uint16
	unitsPerEm         uint16
	created            int64
	modified           int64
	xMin               int16
	yMin               int16
	xMax               int16
	yMax               int16
	macStyle           uint16
	lowestRecPpem      uint16
	fontDirectionHint  int16
	indexToLocFormat   int16
	glyphDataFormat    int16
}

func (t *headTable) tableType() int {
	return Head
}

type glyphTable struct {
	glyphs    []ttfGlyph
	numGlyphs uint16
}

func (t *glyphTable) tableType() int {
	return Glyph
}

type ttfGlyph struct {
	numberOfContours int16
	xMin             int16
	yMin             int16
	xMax             int16
	yMax             int16

	descrip ttfGlyphI

	instructionLength uint16
	instructions      []uint8

	index uint32

	outline *ttfOutline
	bitmap  *ttfBitmap
}

type ttfCompoundGlyph struct {
	comps    []ttfCompoundComp
	numComps uint16
}

func (g *ttfCompoundGlyph) glyphType() int {
	return Compound
}

type ttfCompoundComp struct {
	flags      uint16
	glyphIndex uint16
	arg1       int16
	arg2       int16
	xScale     float32
	yScale     float32
	scale01    float32
	scale10    float32
	xTranslate int16
	yTranslate int16
	point1     uint16
	point2     uint16
}

type ttfSimpleGlyph struct {
	endPtsOfContours  *uint16
	instructionLength uint16
	instructions      []uint8
	flags             *uint8
	xCoordinates      *int16
	yCoordinates      *int16
	numPoints         int16
}

func (g *ttfSimpleGlyph) glyphType() int {
	return Simple
}

type ttfBitmap struct {
	h, w int
	data []uint32
	c    uint32
}

type ttfOutline struct {
	contours    []ttfContour
	numContours int16
	xMin        float32
	yMin        float32
	xMax        float32
	yMax        float32
}

type ttfContour struct {
	segments    []ttfSegment
	numSegments int16
}

type ttfSegment struct {
	segmentType int
	x, y        *float32
	numPOints   int16
}

type fpgmTable struct {
	instructions    []uint8
	numInstructions uint16
}

func (t *fpgmTable) tableType() int {
	return Fpgm
}

type cvtTable struct {
	controlValues []int16
	numValues     uint16
}

func (t *cvtTable) tableType() int {
	return Cvt
}

type cmapTable struct {
	version      uint16
	numSubtables uint16

	subtables []cmapSubTable
}

func (t *cmapTable) tableType() int {
	return Cmap
}

type cmapSubTable struct {
	platformID          uint16
	platformSpecifiedID uint16
	offset              uint32

	format   uint32
	length   uint32
	language uint32

	numIndicies     uint16
	glyphIndexArray []uint32
}

func parseFile(font *ttfFont, filename string) error {
	var err error
	font.fd, err = os.Open(filename)
	check(err)
	defer closeAndNil(&font.fd)

	err = readFontDir(font)
	if err != nil {
		return err
	}
	err = loadTables(font)
	if err != nil {
		return err
	}
	return nil
}

func readFontDir(font *ttfFont) error {
	var seek, err = font.fd.Seek(0, os.SEEK_SET)
	println(seek)
	check(err)
	if seek != 0 {
		panic("failed to seek font dir")
	}

	var reader = font.fd

	font.scalerType = readFixed(reader)
	font.numTables = readUShort(reader)
	font.searchRange = readUShort(reader)
	font.entrySelector = readUShort(reader)
	font.rangeShift = readUShort(reader)

	seek, err = font.fd.Seek(0, os.SEEK_CUR)
	check(err)

	if seek != 12 {
		return fmt.Errorf("incorrect number of bytes in offset subtable: %v", seek)
	}

	font.tables = make([]ttfTable, font.numTables)
	if len(font.tables) == 0 {
		return fmt.Errorf("failed to alloc font tables")
	}

	for i := range font.tables {
		var table = font.tables[i]

		table.tag = readTag(reader)
		table.checkSum = readULong(reader)
		table.offset = readULong(reader)
		table.length = readULong(reader)

		table.status = statusNone
	}

	return nil
}

func loadTables(font *ttfFont) error {
	var requiredTables = []uint32{
		Head,
		Hhea,
		Maxp,
		Post,
		Loca,
	}

	for i := range requiredTables {
		if err := loadTable(font, getTable(font, requiredTables[i])); err != nil {
			return err
		}
	}
	return nil
}

func loadTable(font *ttfFont, table *ttfTable) error {
	var _, err = font.fd.Seek(int64(table.offset), os.SEEK_SET)
	check(err)
	switch table.tag {
	case 0x322f534f: /* OS/2 */
		break
	case 0x544c4350: /* PCLT */
		break
	case 0x70616d63: /* cmap */
		loadCmapTable(font, table)
		break
	case 0x20747663: /* cvt  */
		loadCvtTable(font, table)
		break
	case 0x6d677066: /* fpgm */
		loadFpgmTable(font, table)
		break
	case 0x70736167: /* gasp */
		break
	case 0x66796c67: /* glyph */
		loadGlyphTable(font, table)
		break
	case 0x786d6468: /* hdmx */
		break
	case 0x64616568: /* head */
		// loadHeadTable(font, table)
		break
	case 0x61656868: /* hhea */
		// loadHheaTable(font, table)
		break
	case 0x78746d68: /* hmtx */
		// loadHmtxTable(font, table)
		break
	case 0x6e72656b: /* kern */
		break
	case 0x61636f6c: /* loca */
		// loadLocaTable(font, table)
		break
	case 0x7078616d: /* maxp */
		// loadMaxpTable(font, table)
		break
	case 0x656d616e: /* name */
		break
	case 0x74736f70: /* post */
		// loadPostTable(font, table)
		break
	case 0x70657270: /* prep */
		break
	default:
		return fmt.Errorf("unknown font table type '%v'", table.tag)
		break
	}
	return nil
}

func loadCmapTable(font *ttfFont, table *ttfTable) {
	var cmap = table.data.(*cmapTable)
	var reader = bufio.NewReader(font.fd)

	cmap.version = readUShort(reader)
	cmap.numSubtables = readUShort(reader)
	cmap.subtables = make([]cmapSubTable, cmap.numSubtables)

	var pos int64
	var err error

	for i := range cmap.subtables {
		var subtable = &cmap.subtables[i]

		subtable.platformID = readUShort(reader)
		subtable.platformSpecifiedID = readUShort(reader)
		subtable.offset = readULong(reader)

		pos, err = font.fd.Seek(0, os.SEEK_CUR)
		check(err)

		_, err = font.fd.Seek(int64(table.offset+subtable.offset), os.SEEK_SET)
		check(err)

		loadCmapSubTable(font, subtable)

		_, err = font.fd.Seek(pos, os.SEEK_SET)
		check(err)
	}
}

func loadCmapSubTable(font *ttfFont, subtable *cmapSubTable) {
	var maxp = getMaxpTable(font)
	var reader = bufio.NewReader(font.fd)

	subtable.format = uint32(readUShort(reader))

	if subtable.format < 8 {
		subtable.length = uint32(readUShort(reader))
		subtable.language = uint32(readUShort(reader))
	} else {
		readUShort(reader) // format was actually fixed, this is the .X part
		subtable.length = readULong(reader)
		subtable.language = readULong(reader)
	}

	switch subtable.format {
	case 0:
		subtable.numIndicies = 256
		subtable.glyphIndexArray = make([]uint32, subtable.numIndicies)

		for i := range subtable.glyphIndexArray {
			subtable.glyphIndexArray[i] = uint32((int(readByte(reader)) + 256) % 256)
		}
	case 4:
		var segCount = readUShort(reader)
		var endCode = make([]uint16, segCount)
		var startCode = make([]uint16, segCount)
		var idDelta = make([]uint16, segCount)
		var idRangeOffset = make([]uint16, segCount)

		font.fd.Seek(6, os.SEEK_CUR) //skip 3 words

		var read = func(buf *[]uint16) {
			for i := range *buf {
				(*buf)[i] = readUShort(reader)
			}
		}

		read(&endCode)
		font.fd.Seek(2, os.SEEK_CUR) //skip reserved hword
		read(&startCode)
		read(&idDelta)
		read(&idRangeOffset)

		subtable.numIndicies = maxp.numGlyphs
		subtable.glyphIndexArray = make([]uint32, maxp.numGlyphs)

		var pos, _ = font.fd.Seek(0, os.SEEK_CUR)

		for i := uint16(0); i != segCount; i++ {
			var start = startCode[i]
			var end = endCode[i]
			var delta = idDelta[i]
			var rangeOffset = idRangeOffset[i]

			if start != 65535 && end != 65535 {
				for j := start; j != end; j++ {
					if rangeOffset == 0 {
						subtable.glyphIndexArray[uint32(j+delta)%65536] = uint32(j)
					} else {
						var glyphOffset = pos + ((int64(rangeOffset)/2)+(int64(j)-int64(start))+(int64(i)-int64(segCount)))*2
						font.fd.Seek(glyphOffset, os.SEEK_SET)
						var glyphIndex = readUShort(reader)
						if glyphIndex != 0 {
							glyphIndex = uint16(int(glyphIndex+delta) % 65536)
							if subtable.glyphIndexArray[glyphIndex] == 0 {
								subtable.glyphIndexArray[glyphIndex] = uint32(j)
							}
						}
					}
				}
			}
		}
	default:
		subtable.numIndicies = 0
	}

}

func loadCvtTable(font *ttfFont, table *ttfTable) {
	var cvt = table.data.(*cvtTable)

	cvt.numValues = uint16(table.length / 2)
	cvt.controlValues = make([]int16, cvt.numValues)

	for i := range cvt.controlValues {
		cvt.controlValues[i] = readShort(font.fd)
	}
}

func loadFpgmTable(font *ttfFont, table *ttfTable) {
	var fpgm = table.data.(*fpgmTable)

	fpgm.numInstructions = uint16(table.length)
	fpgm.instructions = make([]uint8, fpgm.numInstructions)

	for i := range fpgm.instructions {
		fpgm.instructions[i] = readByte(font.fd)
	}
}

func loadGlyphTable(font *ttfFont, table *ttfTable) {
	var glyphT = table.data.(*glyphTable)
	var maxp = getMaxpTable(font)
	var loca = getLocaTable(font)

	glyphT.numGlyphs = maxp.numGlyphs
	glyphT.glyphs = make([]ttfGlyph, glyphT.numGlyphs)

	for i := range glyphT.glyphs {
		var glyph = &glyphT.glyphs[i]
		if loca.offsets[i+1]-loca.offsets[i] == 0 {
			continue
		}
		font.fd.Seek(int64(table.offset+loca.offsets[i]), os.SEEK_SET)
		loadGlyph(font, glyph)
	}
}

func loadGlyph(font *ttfFont, glyph *ttfGlyph) {
	glyph.numberOfContours = readShort(font.fd)
	glyph.xMin = readShort(font.fd)
	glyph.yMin = readShort(font.fd)
	glyph.xMax = readShort(font.fd)
	glyph.yMax = readShort(font.fd)

	if glyph.numberOfContours == 0 {
		//Empty glyph
	} else if glyph.numberOfContours > 0 {
		//Simple glyph
		loadSimpleGlyph(font, glyph)
	} else {
		//Compound glyph
		// loadCompoundGlyph(font, glyph)
	}
}

func loadSimpleGlyph(font *ttfFont, glyph *ttfGlyph) {
	// var simpGlyph = glyph.descrip.(*ttfSimpleGlyph)
}

func getTableByName(font *ttfFont, name []uint8) *ttfTable {
	return getTable(font, sToTag(name))
}

func getTable(font *ttfFont, tag uint32) *ttfTable {
	for i := range font.tables {
		if font.tables[i].tag == tag {
			return &font.tables[i]
		}
	}
	return nil
}

func getMaxpTable(font *ttfFont) *maxpTable {
	return getTable(font, Maxp).data.(*maxpTable)
}

func getLocaTable(font *ttfFont) *locaTable {
	return getTable(font, Loca).data.(*locaTable)
}
