package spdx

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/adrg/strutil/metrics"
	"github.com/mitchellh/go-spdx"
)

type syncLicenseList struct {
	*spdx.LicenseList
	*sync.Mutex
}

var Licenses = syncLicenseList{
	LicenseList: &spdx.LicenseList{},
	Mutex:       &sync.Mutex{},
}

func (sll syncLicenseList) License(id string) *spdx.LicenseInfo {
	sll.Lock()
	l := sll.LicenseList.License(id)
	sll.Unlock()
	return l
}

func StartUpdater(ctx context.Context) {
	err := Update()
	if err != nil {
		log.Fatalln("Error updating SPDX license list:", err)
	}

	ticker := time.NewTicker(time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				err = Update()
				if err != nil {
					log.Println("Error updating SPDX license list:", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func Update() error {
	list, err := spdx.List()
	if err != nil {
		return err
	}
	Licenses.Lock()
	Licenses.LicenseList = list
	Licenses.Unlock()

	return nil
}

// findSimilar finds the most similar license ID
// to the one provided
func FindSimilarLicense(s string) string {
	Licenses.Lock()
	defer Licenses.Unlock()

	jw := metrics.NewJaroWinkler()
	jw.CaseSensitive = false

	sims := make([]float64, len(Licenses.Licenses))
	for i, license := range Licenses.Licenses {
		sims[i] = jw.Compare(s, license.ID)
	}

	index := maxIndex(sims)

	if index == -1 {
		return ""
	} else {
		return Licenses.Licenses[index].ID
	}
}

func maxIndex(ff []float64) int {
	if len(ff) == 0 {
		return -1
	}

	m := ff[0]
	mi := 0
	for i, f := range ff {
		if f > m {
			m = f
			mi = i
		}
	}
	return mi
}
