package internal

import (
	"fmt"
	"github.com/0xc0d/vessel/pkg/image"
	"github.com/spf13/cobra"
)

// Images gets all available local images and prints them.
func Images(_ *cobra.Command, _ []string) error {
	imgs, err := image.GetAll()
	if err != nil {
		return err
	}

	pPrintImages(imgs)
	return nil
}

func pPrintImages(imgs []*image.Image) {
	fmt.Println("REPOSITORY\t\t\tTAG\t\tIMAGE ID")
	for _, img := range imgs {
		fmt.Printf("%s\t\t\t%s\t\t%.12s\n", img.Repository, img.Tag, img.ID)
	}
}
