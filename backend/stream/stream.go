package stream

import (
	"io"
	"log"
	"missioncontrol/backend/globals"
)

type Stream interface {
	io.Reader
	io.Writer
	io.Closer
}

// TODO: To optimize allocations of data copying we might be able to use a sync.Pool
func Run(s Stream) error {
	defer close(globals.ByteChannel)
	defer s.Close()

	go func() {
		buf := make([]byte, 256)

		for {
			n, err := s.Read(buf)

			if n > 0 {
				data := make([]byte, n)
				copy(data, buf[:n])

				globals.ByteChannel <- data
			}

			if err != nil {
				if err != io.EOF {
					log.Printf("Stream read error: %v", err)
				}

				return
			}
		}
	}()

	for packetBytes := range globals.RequestChannel {
		_, err := s.Write(packetBytes)

		if err != nil {
			log.Printf("Stream write error: %v", err)
			return err
		}
	}

	return nil
}
