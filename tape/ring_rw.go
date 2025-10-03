package tape

import "os"

func (d Drive) WriteTo(wr *os.File) error {
	ch := make(chan []byte)
	var scsiError error
	go func() {
		for {
			b, err := d.Read()
			if err != nil {
				scsiError = err
				close(ch)
				return
			}
			ch <- b
		}
	}()
	for b := range ch {
		_, err := wr.Write(b)
		if err != nil {
			// TODO cancel
			return err
		}
	}
	return scsiError
}
