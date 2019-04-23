package utils

import (
	"fmt"
	"os"

	"github.com/cloud66/trackman/utils"
	"github.com/khash/updater"
)

func UpdateExec(channel string) {
	update(false, channel)

	fmt.Println("Updated")
}

func update(silent bool, channel string) {
	fmt.Println("Inside the update method with channel : ", channel)
	worker, err := updater.NewUpdater(utils.Version, &updater.Options{
		RemoteURL: "https://s3.amazonaws.com/downloads.cloud66.com/starter/",
		Channel:   utils.Channel,
		Silent:    silent,
	})
	if err != nil {
		if !silent {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	err = worker.Run(channel != utils.Channel)
	if err != nil {
		if !silent {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}
}
