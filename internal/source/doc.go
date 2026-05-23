// Package source implements log source management for watchlog.
//
// It provides a Tailer that can follow one or more log files in real time,
// automatically re-opening them on rotation. When no file paths are supplied
// the Tailer falls back to reading lines from standard input, making it
// easy to pipe log streams directly into watchlog.
//
// Usage:
//
//	tlr, err := source.New([]string{"/var/log/app.log"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer tlr.Close()
//
//	for line := range tlr.Lines() {
//		if line.Err != nil {
//			fmt.Fprintln(os.Stderr, "source error:", line.Err)
//			continue
//		}
//		fmt.Println(line.Source, line.Text)
//	}
package source
