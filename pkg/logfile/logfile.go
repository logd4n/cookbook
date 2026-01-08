/*  # # #   # # #  # # #  # # # #       # # # #     #       #        # # #   # # # #
  #     # #     # #    # #             #    # #   #        #       # * * #  #
 #       #     # #    # # # # #       # # #   # #         #       #  *  #  #  # # #
#     # #     # #    # #             #    #   #          #       # * * #  #      #
# # #   # # #  # # #  # # # #       # # #    #          # # # #  # # #    # # # # */
//yes, log by LOG :)

// ver 1.0

package logfile

import (
	"errors"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger(path string) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.New("Error in package \"logfile\": open file failed! (" + err.Error() + ")")
	}

	return &Logger{file: file}, nil
}

func (l *Logger) Write(data string) error {
	data = time.Now().Format("02.01.2006 15:04:05 MST: ") + data

	if _, err := l.file.WriteString(data); err != nil {
		return errors.New("Error in package \"logfile\": write to file failed! (" + err.Error() + ")")
	}

	return nil
}

func (l *Logger) WriteNewLine() error {
	if _, err := l.file.WriteString("\n"); err != nil {
		return errors.New("Error in package \"logfile\": write to file failed! (" + err.Error() + ")")
	}

	return nil
}

func (l *Logger) Close() {
	l.file.Close()
}
