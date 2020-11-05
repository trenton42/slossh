package slossh

// Recorder receives sessions and passes them on to other integrations
func (s *Slossh) Recorder() {
	for {
		select {
		case session := <-s.recordChan:
			for _, rec := range s.recorders {
				rec.Record(session)
			}
		}
	}
}
