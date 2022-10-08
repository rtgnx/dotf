package main

import (
	"log"
	"os"
	"path"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/blobfs"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/hairyhenderson/go-fsimpl/httpfs"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	godots "github.comr/rtgnx/dotf"
	. "github.comr/rtgnx/dotf/util"
)

var (
	baseDir    string
	plain      bool
	profiles   []string
	installCmd = &cobra.Command{
		Use:   "install [flags] CONFIG VARIABLES",
		Short: "install dot files from source",
		Example: `
      $: dot install --base /home/testuser ./dots ./variables.yml
      $: dot install git+https://github.com/rtgnx/example-dots ./variables.yml
      $: dot install git+https://github.com/rtgnx/example-dots s3://bucket.region.host.com/variables.yml
    `,
		Run:        installCmdRun,
		ArgAliases: []string{"config", "variables"},
		Args:       cobra.MinimumNArgs(2),
	}
)

func init() {
	installCmd.Flags().StringVar(&baseDir, "base", os.Getenv("HOME"), "installation directory")
	installCmd.Flags().StringArrayVar(&profiles, "profiles", []string{}, "profiles to install, default all")
}

func installCmdRun(cmd *cobra.Command, args []string) {
	args[0] = AbsOrScheme(args[0])
	args[1] = AbsOrScheme(args[1])

	dots := godots.Dots([]godots.Dot{})
	vars := &godots.Variables{Global: make(map[string]string)}

	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	mux.Add(httpfs.FS)
	mux.Add(blobfs.FS)
	mux.Add(gitfs.FS)
	log.Println("open: ", args[0])
	configFS := Must(mux.Lookup(args[0]))
	fd := Must(configFS.Open("config.yml"))

	if err := yaml.NewDecoder(Must(configFS.Open("config.yml"))).Decode(&dots); err != nil {
		log.Fatal(err)
	}
	fd.Close()

	log.Println(dots.List())
	varFS := Must(mux.Lookup(path.Dir(args[1])))
	fd = Must(varFS.Open(path.Base(args[1])))

	defer fd.Close()
	if err := vars.ReadIn(fd); err != nil {
		log.Fatal(err)
	}

	if len(profiles) > 0 {
		for _, name := range profiles {
			if err := dots.Run(name, configFS, vars, baseDir); err != nil {
				log.Print(err)
			}
		}
		return
	}

	if err := dots.All(configFS, vars, baseDir); err != nil {
		log.Print(err)
	}

}
