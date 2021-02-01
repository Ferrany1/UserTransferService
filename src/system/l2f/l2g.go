package l2f

import (
	"github.com/Ferrany1/log2file"
	"log"
)

var Log *log.Logger

// Inits a logger file instance
func InitLogger() (Log *log.Logger) {
	Log, err := log2file.GetOptions().Logger()
	if err != nil {
		log.Println(err)
		return nil
	}
	return
}