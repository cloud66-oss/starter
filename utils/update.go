package utils

import (
	"fmt"
	"os"

	"github.com/khash/updater"
)

func UpdateExec(channel string) {
	update(false, channel)

	fmt.Println("Updated")
}

func update(silent bool, channel string) {
	worker, err := updater.NewUpdater(Version, &updater.Options{
		RemoteURL: "https://s3.amazonaws.com/downloads.cloud66.com/starter/",
		Channel:   Channel,
		Silent:    silent,
	})
	if err != nil {
		if !silent {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	err = worker.Run(channel != Channel)
	if err != nil {
		if !silent {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}
}
