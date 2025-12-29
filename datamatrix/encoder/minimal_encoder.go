package encoder

import (
	"fmt"
	"math"
	"strings"

	"github.com/makiuchi-d/gozxing/common"
	"golang.org/x/text/encoding"
)

// Mode represents the encoding mode
type Mode int

const (
	ModeASCII Mode = iota
	ModeC40
	ModeTEXT
	ModeX12
	ModeEDF
	ModeB256
)

var c40Shift2Chars = []rune{'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	':', ';', '<', '=', '>', '?', '@', '[', '\\', ']', '^', '_'}

func isExtendedASCII(ch rune, fnc1 int) bool {
	return ch != rune(fnc1) && ch >= 128 && ch <= 255
}

func isInC40Shift1Set(ch rune) bool {
	return ch <= 31
}

func isInC40Shift2Set(ch rune, fnc1 int) bool {
	for _, c40Shift2Char := range c40Shift2Chars {
		if c40Shift2Char == ch {
			return true
		}
	}
	return ch == rune(fnc1)
}

func isInTextShift1Set(ch rune) bool {
	return isInC40Shift1Set(ch)
}

func isInTextShift2Set(ch rune, fnc1 int) bool {
	return isInC40Shift2Set(ch, fnc1)
}

// EncodeHighLevel Performs message encoding of a DataMatrix message
//
// @param msg the message
// @return the encoded message (the char values range from 0 to 255)
func EncodeHighLevelMinimal(msg string) ([]byte, error) {
	return EncodeHighLevelMinimalWithOptions(msg, nil, -1, SymbolShapeHint_FORCE_NONE)
}

// EncodeHighLevelMinimalWithOptions Performs message encoding of a DataMatrix message
//
// @param msg the message
// @param priorityCharset The preferred encoding.Encoding. When the value of the argument is nil, the algorithm
//
//	chooses charsets that leads to a minimal representation. Otherwise the algorithm will use the priority
//	charset to encode any character in the input that can be encoded by it if the charset is among the
//	supported charsets.
//
// @param fnc1 denotes the character in the input that represents the FNC1 character or -1 if this is not a GS1
//
//	bar code. If the value is not -1 then a FNC1 is also prepended.
//
// @param shape requested shape.
// @return the encoded message (the char values range from 0 to 255)
func EncodeHighLevelMinimalWithOptions(msg string, priorityCharset encoding.Encoding, fnc1 int, shape SymbolShapeHint) ([]byte, error) {
	macroId := 0
	if strings.HasPrefix(msg, HighLevelEncoder_MACRO_05_HEADER) && strings.HasSuffix(msg, HighLevelEncoder_MACRO_TRAILER) {
		macroId = 5
		msg = msg[len(HighLevelEncoder_MACRO_05_HEADER) : len(msg)-2]
	} else if strings.HasPrefix(msg, HighLevelEncoder_MACRO_06_HEADER) && strings.HasSuffix(msg, HighLevelEncoder_MACRO_TRAILER) {
		macroId = 6
		msg = msg[len(HighLevelEncoder_MACRO_06_HEADER) : len(msg)-2]
	}
	return encodeMinimal(msg, priorityCharset, fnc1, shape, macroId), nil
}

// encodeMinimal Encodes input minimally and returns an array of the codewords
//
// @param input The string to encode
// @param priorityCharset The preferred encoding.Encoding. When the value of the argument is nil, the algorithm
//
//	chooses charsets that leads to a minimal representation. Otherwise the algorithm will use the priority
//	charset to encode any character in the input that can be encoded by it if the charset is among the
//	supported charsets.
//
// @param fnc1 denotes the character in the input that represents the FNC1 character or -1 if this is not a GS1
//
//	bar code. If the value is not -1 then a FNC1 is also prepended.
//
// @param shape requested shape.
// @param macroId Prepends the specified macro function in case that a value of 5 or 6 is specified.
// @return An array of bytes representing the codewords of a minimal encoding.
func encodeMinimal(input string, priorityCharset encoding.Encoding, fnc1 int, shape SymbolShapeHint, macroId int) []byte {
	return encodeMinimally(newInput(input, priorityCharset, fnc1, shape, macroId)).getBytes()
}

func addEdge(edges [][]*edge, edge *edge) {
	vertexIndex := edge.fromPosition + edge.characterLength
	endModeOrdinal := int(edge.getEndMode())
	if edges[vertexIndex][endModeOrdinal] == nil ||
		edges[vertexIndex][endModeOrdinal].cachedTotalSize > edge.cachedTotalSize {
		edges[vertexIndex][endModeOrdinal] = edge
	}
}

// getNumberOfC40Words returns the number of words in which the string starting at from can be encoded in c40 or text mode.
// The number of characters encoded is returned in characterLength.
// The number of characters encoded is also minimal in the sense that the algorithm stops as soon
// as a character encoding fills a C40 word competely (three C40 values). An exception is at the
// end of the string where two C40 values are allowed (according to the spec the third c40 value
// is filled  with 0 (Shift 1) in this case).
func getNumberOfC40Words(input *input, from int, c40 bool, characterLength *int) int {
	thirdsCount := 0
	for i := from; i < input.length(); i++ {
		if isECI, _ := input.isECI(i); isECI {
			*characterLength = 0
			return 0
		}
		ci, _ := input.charAt(i)
		if (c40 && isNativeC40Char(ci)) || (!c40 && isNativeTextChar(ci)) {
			thirdsCount++ // native
		} else if !isExtendedASCII(ci, input.getFNC1Character()) {
			thirdsCount += 2 // shift
		} else {
			asciiValue := int(ci) & 0xff
			if asciiValue >= 128 && ((c40 && isNativeC40Char(rune(asciiValue-128))) ||
				(!c40 && isNativeTextChar(rune(asciiValue-128)))) {
				thirdsCount += 3 // shift, Upper shift
			} else {
				thirdsCount += 4 // shift, Upper shift, shift
			}
		}

		if thirdsCount%3 == 0 || ((thirdsCount-2)%3 == 0 && i+1 == input.length()) {
			*characterLength = i - from + 1
			return int(math.Ceil(float64(thirdsCount) / 3.0))
		}
	}
	*characterLength = 0
	return 0
}

func isNativeC40Char(ch rune) bool {
	return (ch == ' ') || (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z')
}

func isNativeTextChar(ch rune) bool {
	return (ch == ' ') || (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z')
}

func isNativeX12Char(ch rune) bool {
	return isX12TermSepChar(ch) || (ch == ' ') || (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z')
}

func isX12TermSepChar(ch rune) bool {
	return (ch == '\r') || // CR
		(ch == '*') ||
		(ch == '>')
}

func isNativeEDIFACTChar(ch rune) bool {
	return ch >= ' ' && ch <= '^'
}

func isDigitChar(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func addEdges(input *input, edges [][]*edge, from int, previous *edge) {
	if isECI, _ := input.isECI(from); isECI {
		addEdge(edges, newEdge(input, ModeASCII, from, 1, previous))
		return
	}

	ch, _ := input.charAt(from)
	if previous == nil || previous.getEndMode() != ModeEDF { // not possible to unlatch a full EDF edge to something else
		if isDigitChar(ch) && input.haveNCharacters(from, 2) {
			ch2, _ := input.charAt(from + 1)
			if isDigitChar(ch2) {
				// two digits ASCII encoded
				addEdge(edges, newEdge(input, ModeASCII, from, 2, previous))
			} else {
				// one ASCII encoded character or an extended character via Upper Shift
				addEdge(edges, newEdge(input, ModeASCII, from, 1, previous))
			}
		} else {
			// one ASCII encoded character or an extended character via Upper Shift
			addEdge(edges, newEdge(input, ModeASCII, from, 1, previous))
		}

		modes := []Mode{ModeC40, ModeTEXT}
		for _, mode := range modes {
			var characterLength int
			if getNumberOfC40Words(input, from, mode == ModeC40, &characterLength) > 0 {
				addEdge(edges, newEdge(input, mode, from, characterLength, previous))
			}
		}

		if input.haveNCharacters(from, 3) {
			ch1, _ := input.charAt(from)
			ch2, _ := input.charAt(from + 1)
			ch3, _ := input.charAt(from + 2)
			if isNativeX12Char(ch1) && isNativeX12Char(ch2) && isNativeX12Char(ch3) {
				addEdge(edges, newEdge(input, ModeX12, from, 3, previous))
			}
		}

		addEdge(edges, newEdge(input, ModeB256, from, 1, previous))
	}

	// We create 4 EDF edges, with 1, 2 3 or 4 characters length. The fourth normally doesn't have a latch to ASCII
	// unless it is 2 characters away from the end of the input.
	var i int
	for i = 0; i < 3; i++ {
		pos := from + i
		if input.haveNCharacters(pos, 1) {
			ch, _ := input.charAt(pos)
			if isNativeEDIFACTChar(ch) {
				addEdge(edges, newEdge(input, ModeEDF, from, i+1, previous))
			} else {
				break
			}
		} else {
			break
		}
	}
	if i == 3 && input.haveNCharacters(from, 4) {
		ch, _ := input.charAt(from + 3)
		if isNativeEDIFACTChar(ch) {
			addEdge(edges, newEdge(input, ModeEDF, from, 4, previous))
		}
	}
}

func encodeMinimally(input *input) *result {
	inputLength := input.length()

	// Array that represents vertices. There is a vertex for every character and mode.
	// The last dimension in the array below encodes the 6 modes ASCII, C40, TEXT, X12, EDF and B256
	edges := make([][]*edge, inputLength+1)
	for i := range edges {
		edges[i] = make([]*edge, 6)
	}
	addEdges(input, edges, 0, nil)

	for i := 1; i <= inputLength; i++ {
		for j := 0; j < 6; j++ {
			if edges[i][j] != nil && i < inputLength {
				addEdges(input, edges, i, edges[i][j])
			}
		}
		// optimize memory by removing edges that have been passed.
		for j := 0; j < 6; j++ {
			edges[i-1][j] = nil
		}
	}

	minimalJ := -1
	minimalSize := int(^uint(0) >> 1) // MaxInt
	for j := 0; j < 6; j++ {
		if edges[inputLength][j] != nil {
			edge := edges[inputLength][j]
			size := edge.cachedTotalSize
			if j >= 1 && j <= 3 { // C40, TEXT and X12 need an extra unlatch at the end
				size++
			}
			if size < minimalSize {
				minimalSize = size
				minimalJ = j
			}
		}
	}

	if minimalJ < 0 {
		panic(fmt.Sprintf("Failed to encode \"%s\"", input.String()))
	}
	return newResult(edges[inputLength][minimalJ])
}

type edge struct {
	input           *input
	mode            Mode // the mode at the start of this edge.
	fromPosition    int
	characterLength int
	previous        *edge
	cachedTotalSize int
}

var (
	allCodewordCapacities         = []int{3, 5, 8, 10, 12, 16, 18, 22, 30, 32, 36, 44, 49, 62, 86, 114, 144, 174, 204, 280, 368, 456, 576, 696, 816, 1050, 1304, 1558}
	squareCodewordCapacities      = []int{3, 5, 8, 12, 18, 22, 30, 36, 44, 62, 86, 114, 144, 174, 204, 280, 368, 456, 576, 696, 816, 1050, 1304, 1558}
	rectangularCodewordCapacities = []int{5, 10, 16, 33, 32, 49}
)

func newEdge(input *input, mode Mode, fromPosition int, characterLength int, previous *edge) *edge {
	size := 0
	if previous != nil {
		size = previous.cachedTotalSize
	}

	previousMode := getPreviousMode(previous)

	// Switching modes
	// ASCII -> C40: latch 230
	// ASCII -> TEXT: latch 239
	// ASCII -> X12: latch 238
	// ASCII -> EDF: latch 240
	// ASCII -> B256: latch 231
	// C40 -> ASCII: word(c1,c2,c3), 254
	// TEXT -> ASCII: word(c1,c2,c3), 254
	// X12 -> ASCII: word(c1,c2,c3), 254
	// EDIFACT -> ASCII: Unlatch character,0,0,0 or c1,Unlatch character,0,0 or c1,c2,Unlatch character,0 or
	// c1,c2,c3,Unlatch character
	// B256 -> ASCII: without latch after n bytes
	switch mode {
	case ModeASCII:
		size++
		if isECI, _ := input.isECI(fromPosition); isECI {
			size++
		} else {
			ch, _ := input.charAt(fromPosition)
			if isExtendedASCII(ch, input.getFNC1Character()) {
				size++
			}
		}
		if previousMode == ModeC40 || previousMode == ModeTEXT || previousMode == ModeX12 {
			size++ // unlatch 254 to ASCII
		}
	case ModeB256:
		size++
		if previousMode != ModeB256 {
			size++ // byte count
		} else if getB256Size(previous) == 250 {
			size++ // extra byte count
		}
		if previousMode == ModeASCII {
			size++ // latch to B256
		} else if previousMode == ModeC40 || previousMode == ModeTEXT || previousMode == ModeX12 {
			size += 2 // unlatch to ASCII, latch to B256
		}
	case ModeC40, ModeTEXT, ModeX12:
		if mode == ModeX12 {
			size += 2
		} else {
			var charLen int
			size += getNumberOfC40Words(input, fromPosition, mode == ModeC40, &charLen) * 2
		}

		if previousMode == ModeASCII || previousMode == ModeB256 {
			size++ // additional byte for latch from ASCII to this mode
		} else if previousMode != mode && (previousMode == ModeC40 || previousMode == ModeTEXT || previousMode == ModeX12) {
			size += 2 // unlatch 254 to ASCII followed by latch to this mode
		}
	case ModeEDF:
		size += 3
		if previousMode == ModeASCII || previousMode == ModeB256 {
			size++ // additional byte for latch from ASCII to this mode
		} else if previousMode == ModeC40 || previousMode == ModeTEXT || previousMode == ModeX12 {
			size += 2 // unlatch 254 to ASCII followed by latch to this mode
		}
	}

	return &edge{
		input:           input,
		mode:            mode,
		fromPosition:    fromPosition,
		characterLength: characterLength,
		previous:        previous,
		cachedTotalSize: size,
	}
}

// getB256Size does not count beyond 250
func getB256Size(e *edge) int {
	cnt := 0
	current := e
	for current != nil && current.mode == ModeB256 && cnt <= 250 {
		cnt++
		current = current.previous
	}
	return cnt
}

func getPreviousStartMode(e *edge) Mode {
	if e == nil {
		return ModeASCII
	}
	return e.mode
}

func getPreviousMode(e *edge) Mode {
	if e == nil {
		return ModeASCII
	}
	return e.getEndMode()
}

// getEndMode Returns ModeASCII in case that:
//   - Mode is EDIFACT and characterLength is less than 4 or the remaining characters can be encoded in at most 2
//     ASCII bytes.
//   - Mode is C40, TEXT or X12 and the remaining characters can be encoded in at most 1 ASCII byte.
//     Returns mode in all other cases.
func (e *edge) getEndMode() Mode {
	if e.mode == ModeEDF {
		if e.characterLength < 4 {
			return ModeASCII
		}
		lastASCII := e.getLastASCII() // see 5.2.8.2 EDIFACT encodation Rules
		if lastASCII > 0 && e.getCodewordsRemaining(e.cachedTotalSize+lastASCII) <= 2-lastASCII {
			return ModeASCII
		}
	}
	if e.mode == ModeC40 || e.mode == ModeTEXT || e.mode == ModeX12 {
		// see 5.2.5.2 C40 encodation rules and 5.2.7.2 ANSI X12 encodation rules
		if e.fromPosition+e.characterLength >= e.input.length() && e.getCodewordsRemaining(e.cachedTotalSize) == 0 {
			return ModeASCII
		}
		lastASCII := e.getLastASCII()
		if lastASCII == 1 && e.getCodewordsRemaining(e.cachedTotalSize+1) == 0 {
			return ModeASCII
		}
	}
	return e.mode
}

func (e *edge) getMode() Mode {
	return e.mode
}

// getLastASCII Peeks ahead and returns 1 if the postfix consists of exactly two digits, 2 if the postfix consists of exactly
// two consecutive digits and a non extended character or of 4 digits.
// Returns 0 in any other case
func (e *edge) getLastASCII() int {
	length := e.input.length()
	from := e.fromPosition + e.characterLength
	if length-from > 4 || from >= length {
		return 0
	}
	if length-from == 1 {
		ch, _ := e.input.charAt(from)
		if isExtendedASCII(ch, e.input.getFNC1Character()) {
			return 0
		}
		return 1
	}
	if length-from == 2 {
		ch1, _ := e.input.charAt(from)
		ch2, _ := e.input.charAt(from + 1)
		if isExtendedASCII(ch1, e.input.getFNC1Character()) || isExtendedASCII(ch2, e.input.getFNC1Character()) {
			return 0
		}
		if isDigitChar(ch1) && isDigitChar(ch2) {
			return 1
		}
		return 2
	}
	if length-from == 3 {
		ch1, _ := e.input.charAt(from)
		ch2, _ := e.input.charAt(from + 1)
		ch3, _ := e.input.charAt(from + 2)
		if isDigitChar(ch1) && isDigitChar(ch2) && !isExtendedASCII(ch3, e.input.getFNC1Character()) {
			return 2
		}
		if isDigitChar(ch2) && isDigitChar(ch3) && !isExtendedASCII(ch1, e.input.getFNC1Character()) {
			return 2
		}
		return 0
	}
	ch1, _ := e.input.charAt(from)
	ch2, _ := e.input.charAt(from + 1)
	ch3, _ := e.input.charAt(from + 2)
	ch4, _ := e.input.charAt(from + 3)
	if isDigitChar(ch1) && isDigitChar(ch2) && isDigitChar(ch3) && isDigitChar(ch4) {
		return 2
	}
	return 0
}

// getMinSymbolSize Returns the capacity in codewords of the smallest symbol that has enough capacity to fit the given minimal
// number of codewords.
func (e *edge) getMinSymbolSize(minimum int) int {
	switch e.input.getShapeHint() {
	case SymbolShapeHint_FORCE_SQUARE:
		for _, capacity := range squareCodewordCapacities {
			if capacity >= minimum {
				return capacity
			}
		}
	case SymbolShapeHint_FORCE_RECTANGLE:
		for _, capacity := range rectangularCodewordCapacities {
			if capacity >= minimum {
				return capacity
			}
		}
	}
	for _, capacity := range allCodewordCapacities {
		if capacity >= minimum {
			return capacity
		}
	}
	return allCodewordCapacities[len(allCodewordCapacities)-1]
}

// getCodewordsRemaining Returns the remaining capacity in codewords of the smallest symbol that has enough capacity to fit the given
// minimal number of codewords.
func (e *edge) getCodewordsRemaining(minimum int) int {
	return e.getMinSymbolSize(minimum) - minimum
}

func edgeGetBytes(c int) []byte {
	return []byte{byte(c)}
}

func edgeGetBytes2(c1, c2 int) []byte {
	return []byte{byte(c1), byte(c2)}
}

func setC40Word(bytes []byte, offset int, c1, c2, c3 int) {
	val16 := (1600 * (c1 & 0xff)) + (40 * (c2 & 0xff)) + (c3 & 0xff) + 1
	bytes[offset] = byte(val16 / 256)
	bytes[offset+1] = byte(val16 % 256)
}

func getX12Value(c rune) int {
	if c == 13 {
		return 0
	}
	if c == 42 {
		return 1
	}
	if c == 62 {
		return 2
	}
	if c == 32 {
		return 3
	}
	if c >= 48 && c <= 57 {
		return int(c) - 44
	}
	if c >= 65 && c <= 90 {
		return int(c) - 51
	}
	return int(c)
}

func (e *edge) getX12Words() []byte {
	result := make([]byte, e.characterLength/3*2)
	for i := 0; i < len(result); i += 2 {
		ch1, _ := e.input.charAt(e.fromPosition + i/2*3)
		ch2, _ := e.input.charAt(e.fromPosition + i/2*3 + 1)
		ch3, _ := e.input.charAt(e.fromPosition + i/2*3 + 2)
		setC40Word(result, i, getX12Value(ch1), getX12Value(ch2), getX12Value(ch3))
	}
	return result
}

func getShiftValue(c rune, c40 bool, fnc1 int) int {
	if (c40 && isInC40Shift1Set(c)) || (!c40 && isInTextShift1Set(c)) {
		return 0
	}
	if (c40 && isInC40Shift2Set(c, fnc1)) || (!c40 && isInTextShift2Set(c, fnc1)) {
		return 1
	}
	return 2
}

func getC40Value(c40 bool, setIndex int, c rune, fnc1 int) int {
	if c == rune(fnc1) {
		return 27
	}
	if c40 {
		if c <= 31 {
			return int(c)
		}
		if c == 32 {
			return 3
		}
		if c <= 47 {
			return int(c) - 33
		}
		if c <= 57 {
			return int(c) - 44
		}
		if c <= 64 {
			return int(c) - 43
		}
		if c <= 90 {
			return int(c) - 51
		}
		if c <= 95 {
			return int(c) - 69
		}
		if c <= 127 {
			return int(c) - 96
		}
		return int(c)
	} else {
		if c == 0 {
			return 0
		}
		if setIndex == 0 && c <= 3 {
			return int(c) - 1 // is this a bug in the spec?
		}
		if setIndex == 1 && c <= 31 {
			return int(c)
		}
		if c == 32 {
			return 3
		}
		if c >= 33 && c <= 47 {
			return int(c) - 33
		}
		if c >= 48 && c <= 57 {
			return int(c) - 44
		}
		if c >= 58 && c <= 64 {
			return int(c) - 43
		}
		if c >= 65 && c <= 90 {
			return int(c) - 64
		}
		if c >= 91 && c <= 95 {
			return int(c) - 69
		}
		if c == 96 {
			return 0
		}
		if c >= 97 && c <= 122 {
			return int(c) - 83
		}
		if c >= 123 && c <= 127 {
			return int(c) - 96
		}
		return int(c)
	}
}

func (e *edge) getC40Words(c40 bool, fnc1 int) []byte {
	c40Values := make([]byte, 0)
	for i := 0; i < e.characterLength; i++ {
		ci, _ := e.input.charAt(e.fromPosition + i)
		if (c40 && isNativeC40Char(ci)) || (!c40 && isNativeTextChar(ci)) {
			c40Values = append(c40Values, byte(getC40Value(c40, 0, ci, fnc1)))
		} else if !isExtendedASCII(ci, fnc1) {
			shiftValue := getShiftValue(ci, c40, fnc1)
			c40Values = append(c40Values, byte(shiftValue)) // Shift[123]
			c40Values = append(c40Values, byte(getC40Value(c40, shiftValue, ci, fnc1)))
		} else {
			asciiValue := rune((int(ci) & 0xff) - 128)
			if (c40 && isNativeC40Char(asciiValue)) || (!c40 && isNativeTextChar(asciiValue)) {
				c40Values = append(c40Values, 1)  // Shift 2
				c40Values = append(c40Values, 30) // Upper Shift
				c40Values = append(c40Values, byte(getC40Value(c40, 0, asciiValue, fnc1)))
			} else {
				c40Values = append(c40Values, 1)  // Shift 2
				c40Values = append(c40Values, 30) // Upper Shift
				shiftValue := getShiftValue(asciiValue, c40, fnc1)
				c40Values = append(c40Values, byte(shiftValue)) // Shift[123]
				c40Values = append(c40Values, byte(getC40Value(c40, shiftValue, asciiValue, fnc1)))
			}
		}
	}

	if len(c40Values)%3 != 0 {
		// pad with 0 (Shift 1)
		c40Values = append(c40Values, 0)
	}

	result := make([]byte, len(c40Values)/3*2)
	byteIndex := 0
	for i := 0; i < len(c40Values); i += 3 {
		setC40Word(result, byteIndex, int(c40Values[i]&0xff), int(c40Values[i+1]&0xff), int(c40Values[i+2]&0xff))
		byteIndex += 2
	}
	return result
}

func (e *edge) getEDFBytes() []byte {
	numberOfThirds := int(math.Ceil(float64(e.characterLength) / 4.0))
	result := make([]byte, numberOfThirds*3)
	pos := e.fromPosition
	endPos := e.fromPosition + e.characterLength - 1
	if endPos > e.input.length()-1 {
		endPos = e.input.length() - 1
	}
	for i := 0; i < numberOfThirds; i++ {
		edfValues := make([]int, 4)
		for j := 0; j < 4; j++ {
			if pos <= endPos {
				ch, _ := e.input.charAt(pos)
				edfValues[j] = int(ch) & 0x3f
				pos++
			} else {
				if pos == endPos+1 {
					edfValues[j] = 0x1f
				} else {
					edfValues[j] = 0
				}
			}
		}
		val24 := edfValues[0]<<18 | edfValues[1]<<12 | edfValues[2]<<6 | edfValues[3]
		result[i*3] = byte((val24 >> 16) & 0xff)
		result[i*3+1] = byte((val24 >> 8) & 0xff)
		result[i*3+2] = byte(val24 & 0xff)
	}
	return result
}

func (e *edge) getLatchBytes() []byte {
	previousMode := getPreviousMode(e.previous)
	switch previousMode {
	case ModeASCII, ModeB256: // after B256 ends (via length) we are back to ASCII
		switch e.mode {
		case ModeB256:
			return edgeGetBytes(231)
		case ModeC40:
			return edgeGetBytes(230)
		case ModeTEXT:
			return edgeGetBytes(239)
		case ModeX12:
			return edgeGetBytes(238)
		case ModeEDF:
			return edgeGetBytes(240)
		}
	case ModeC40, ModeTEXT, ModeX12:
		if e.mode != previousMode {
			switch e.mode {
			case ModeASCII:
				return edgeGetBytes(254)
			case ModeB256:
				return edgeGetBytes2(254, 231)
			case ModeC40:
				return edgeGetBytes2(254, 230)
			case ModeTEXT:
				return edgeGetBytes2(254, 239)
			case ModeX12:
				return edgeGetBytes2(254, 238)
			case ModeEDF:
				return edgeGetBytes2(254, 240)
			}
		}
	case ModeEDF:
		// The rightmost EDIFACT edge always contains an unlatch character
		break
	}
	return []byte{}
}

// getDataBytes Important: The function does not return the length bytes (one or two) in case of B256 encoding
func (e *edge) getDataBytes() []byte {
	switch e.mode {
	case ModeASCII:
		if isECI, _ := e.input.isECI(e.fromPosition); isECI {
			eciValue, _ := e.input.getECIValue(e.fromPosition)
			return edgeGetBytes2(241, eciValue+1)
		}
		ch, _ := e.input.charAt(e.fromPosition)
		if isExtendedASCII(ch, e.input.getFNC1Character()) {
			return edgeGetBytes2(235, int(ch)-127)
		}
		if e.characterLength == 2 {
			ch1, _ := e.input.charAt(e.fromPosition)
			ch2, _ := e.input.charAt(e.fromPosition + 1)
			return edgeGetBytes((int(ch1)-'0')*10 + int(ch2) - '0' + 130)
		}
		if isFNC1, _ := e.input.isFNC1(e.fromPosition); isFNC1 {
			return edgeGetBytes(232)
		}
		ch, _ = e.input.charAt(e.fromPosition)
		return edgeGetBytes(int(ch) + 1)
	case ModeB256:
		ch, _ := e.input.charAt(e.fromPosition)
		return edgeGetBytes(int(ch))
	case ModeC40:
		return e.getC40Words(true, e.input.getFNC1Character())
	case ModeTEXT:
		return e.getC40Words(false, e.input.getFNC1Character())
	case ModeX12:
		return e.getX12Words()
	case ModeEDF:
		return e.getEDFBytes()
	}
	return []byte{}
}

type result struct {
	bytes []byte
}

func newResult(solution *edge) *result {
	input := solution.input
	size := 0
	bytesAL := make([]byte, 0)
	randomizePostfixLength := make([]int, 0)
	randomizeLengths := make([]int, 0)
	if (solution.mode == ModeC40 || solution.mode == ModeTEXT || solution.mode == ModeX12) &&
		solution.getEndMode() != ModeASCII {
		size += prepend(edgeGetBytes(254), &bytesAL)
	}
	current := solution
	for current != nil {
		size += prepend(current.getDataBytes(), &bytesAL)

		if current.previous == nil || getPreviousStartMode(current.previous) != current.getMode() {
			if current.getMode() == ModeB256 {
				if size <= 249 {
					bytesAL = append([]byte{byte(size)}, bytesAL...)
					size++
				} else {
					bytesAL = append([]byte{byte(size % 250), byte(size/250 + 249)}, bytesAL...)
					size += 2
				}
				randomizePostfixLength = append(randomizePostfixLength, len(bytesAL))
				randomizeLengths = append(randomizeLengths, size)
			}
			prepend(current.getLatchBytes(), &bytesAL)
			size = 0
		}

		current = current.previous
	}
	if input.getMacroId() == 5 {
		size += prepend(edgeGetBytes(236), &bytesAL)
	} else if input.getMacroId() == 6 {
		size += prepend(edgeGetBytes(237), &bytesAL)
	}

	if input.getFNC1Character() > 0 {
		size += prepend(edgeGetBytes(232), &bytesAL)
	}
	for i := 0; i < len(randomizePostfixLength); i++ {
		applyRandomPattern(&bytesAL, len(bytesAL)-randomizePostfixLength[i], randomizeLengths[i])
	}
	// add padding
	capacity := solution.getMinSymbolSize(len(bytesAL))
	if len(bytesAL) < capacity {
		bytesAL = append(bytesAL, 129)
	}
	for len(bytesAL) < capacity {
		bytesAL = append(bytesAL, byte(randomize253State(len(bytesAL)+1)))
	}

	return &result{bytes: bytesAL}
}

func prepend(bytes []byte, into *[]byte) int {
	for i := len(bytes) - 1; i >= 0; i-- {
		*into = append([]byte{bytes[i]}, *into...)
	}
	return len(bytes)
}

func applyRandomPattern(bytesAL *[]byte, startPosition int, length int) {
	for i := 0; i < length; i++ {
		// See "B.1 253-state algorithm
		padCodewordPosition := startPosition + i
		padCodewordValue := int((*bytesAL)[padCodewordPosition] & 0xff)
		pseudoRandomNumber := ((149 * (padCodewordPosition + 1)) % 255) + 1
		tempVariable := padCodewordValue + pseudoRandomNumber
		if tempVariable <= 255 {
			(*bytesAL)[padCodewordPosition] = byte(tempVariable)
		} else {
			(*bytesAL)[padCodewordPosition] = byte(tempVariable - 256)
		}
	}
}

func (r *result) getBytes() []byte {
	return r.bytes
}

type input struct {
	*common.MinimalECIInput
	shape   SymbolShapeHint
	macroId int
}

func newInput(stringToEncode string, priorityCharset encoding.Encoding, fnc1 int, shape SymbolShapeHint, macroId int) *input {
	return &input{
		MinimalECIInput: common.NewMinimalECIInput(stringToEncode, priorityCharset, fnc1),
		shape:           shape,
		macroId:         macroId,
	}
}

func (in *input) length() int {
	return in.MinimalECIInput.Length()
}

func (in *input) charAt(index int) (rune, error) {
	return in.MinimalECIInput.CharAt(index)
}

func (in *input) isECI(index int) (bool, error) {
	return in.MinimalECIInput.IsECI(index)
}

func (in *input) isFNC1(index int) (bool, error) {
	return in.MinimalECIInput.IsFNC1(index)
}

func (in *input) getECIValue(index int) (int, error) {
	return in.MinimalECIInput.GetECIValue(index)
}

func (in *input) haveNCharacters(index, n int) bool {
	return in.MinimalECIInput.HaveNCharacters(index, n)
}

func (in *input) getFNC1Character() int {
	return in.MinimalECIInput.GetFNC1Character()
}

func (in *input) getMacroId() int {
	return in.macroId
}

func (in *input) getShapeHint() SymbolShapeHint {
	return in.shape
}

func (in *input) String() string {
	return fmt.Sprintf("Input{length=%d}", in.length())
}
