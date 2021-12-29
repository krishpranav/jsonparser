package jsonparser

type scratch struct {
	data []byte
	fill int
}

func (s *scratch) reset() {
	s.fill = 0
}

func (s *scratch) bytes() []byte {
	return s.data[0:s.fill]
}

func (s *scratch) grow() {
	ndata := make([]byte, cap(s.data)*2)
	copy(ndata, s.data[:])
	s.data = ndata
}

func (s *scratch) add(c byte) {
	if s.fill+1 >= cap(s.data) {
		s.grow()
	}

	s.data[s.fill] = c
	s.fill++
}

func (s *scratch) addRune(r rune) int {

}
