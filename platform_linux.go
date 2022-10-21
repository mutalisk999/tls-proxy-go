package tls_proxy_go

import (
	"log"
	"syscall"
)

func SetRLimit(v uint64) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatalln("Error Getting Rlimit:", err)
	}
	log.Printf("Rlimit Current: %d", rLimit.Cur)

	if rLimit.Cur >= v {
		return
	} else {
		rLimit.Cur = v
		rLimit.Max = 999999

		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			log.Fatalln("Error Setting Rlimit:", err)
		}
		log.Printf("Setting Rlimit: %d", rLimit.Cur)

		err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			log.Fatalln("Error Getting Rlimit:", err)
		}
		log.Printf("Rlimit Final: %d", rLimit.Cur)
	}
}
