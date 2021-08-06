package wavreaderseeker

// Object to satisfy ebiten's requirements for an io.Reader for sound.
// This all relies on the wav files being the correct (16 bit stereo) format.

import (
	"errors"
	"io"
)

const DATA_START_POS = 44

type WavReaderSeeker struct {
	Data []byte
	Pos int
}

func NewWavReaderSeeker(bytes []byte) *WavReaderSeeker {
	return &WavReaderSeeker{Data: bytes, Pos: DATA_START_POS}
}

func (self *WavReaderSeeker) Read(p []byte) (n int, err error) {

	// Read reads up to len(p) bytes into p.
	// It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.

	if self.Pos >= len(self.Data) {
		return 0, io.EOF
	}

	advance := copy(p, self.Data[self.Pos:])		// Copy returns the number of elements copied, which will be the minimum of len(dst) and len(src)
	self.Pos += advance

	return advance, nil
}

func (self *WavReaderSeeker) Seek(offset int64, whence int) (int64, error) {

	// Seek sets the offset for the next Read or Write to offset, interpreted according to whence.
	// Seek returns the new offset relative to the start of the file and an error, if any.
	//
	// Seeking to an offset before the start of the file is an error.
	// Seeking to any positive offset is legal, but the behavior of subsequent I/O operations is implementation-dependent.

	var new_pos int

	switch whence {
	case io.SeekStart:
		new_pos = DATA_START_POS
	case io.SeekCurrent:
		new_pos = self.Pos
	case io.SeekEnd:
		new_pos = len(self.Data)
	default:
		panic("Invalid whence")
	}

	new_pos += int(offset)

	if new_pos < DATA_START_POS {
		return 0, errors.New("Attempt to seek to before the data start")
	}

	self.Pos = new_pos

	return int64(self.Pos), nil
}
