package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go_tools/files/compress/gzip"
	copy2 "go_tools/files/copy"
	"go_tools/log"
	"strconv"
	"time"
)

var copyCmd = &cobra.Command{
	Use:     fmt.Sprintf(CUT_OFF_COMMAND_PREFIX, "copy"),
	Aliases: []string{fmt.Sprintf(CUT_OFF_COMMAND_PREFIX, "copy")},
	Short:   "copy files command, speed up copy process by multiple parallelism",
	Long:    "copy files command, speed up copy process by multiple parallelism",
	Example: "co_copy ./aa/ ./bb/ -parallelism 4",
	//Args:                       cobra.ExactArgs(),
	ArgAliases: nil,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := cmd.Flag("s").Value.String()
		targetPath := cmd.Flag("t").Value.String()
		handler, err := copy2.Get(sourcePath, targetPath)
		if err != nil {
			log.Warn("init copy env error: %v", err)
			return
		}
		log.Info("source path: %v, target path: %v", sourcePath, targetPath)
		startTime := time.Now().Unix()
		_, err = handler.Copy()
		if err != nil {
			log.Warn("copy file error: %v", err)
			return
		}
		log.Info("copy process finish, dir: %v, files: %v cost time: %vs", handler.DirNumber, handler.FileNumber, time.Now().Unix()-startTime)
	},
	RunE:                       nil,
	PostRun:                    nil,
	PostRunE:                   nil,
	PersistentPostRun:          nil,
	PersistentPostRunE:         nil,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	CompletionOptions:          cobra.CompletionOptions{},
	TraverseChildren:           false,
	Hidden:                     false,
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
}

var compressCmd = &cobra.Command{
	Use:     fmt.Sprintf(CUT_OFF_COMMAND_PREFIX, "compress"),
	Aliases: []string{fmt.Sprintf(CUT_OFF_COMMAND_PREFIX, "compress")},
	Short:   "compress files command, speed up compress process by multiple parallelism",
	Long:    "compress files command, speed up compress process by multiple parallelism",
	Example: "co_compress ./aa/ ./bb/ -parallelism 4",
	//Args:                       cobra.ExactArgs(),
	ArgAliases: nil,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := cmd.Flag("s").Value.String()
		targetPath := cmd.Flag("t").Value.String()
		parallelism := cmd.Flag("p").Value.String()
		var parallelismNumber int = 0
		var err error
		if len(parallelism) > 0 {
			parallelismNumber, err = strconv.Atoi(parallelism)
			if err != nil {
				panic(fmt.Sprintf("decode flag --p error: %v", err))
			}
		}
		handler, err := gzip.Get(sourcePath, targetPath, parallelismNumber, 0, true, true)
		if err != nil {
			log.Warn("init compress env error: %v", err)
			return
		}
		log.Info("source path: %v, target path: %v, parallelism: %v", sourcePath, targetPath, handler.Parallelism)
		startTime := time.Now().Unix()
		err = handler.Compress()
		if err != nil {
			log.Warn("compress file error: %v", err)
			return
		}
		log.Info("compress process finish, dir: %v, files: %v cost time: %vs", handler.DirNum, handler.FileNum, time.Now().Unix()-startTime)
	},
	RunE:                       nil,
	PostRun:                    nil,
	PostRunE:                   nil,
	PersistentPostRun:          nil,
	PersistentPostRunE:         nil,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	CompletionOptions:          cobra.CompletionOptions{},
	TraverseChildren:           false,
	Hidden:                     false,
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
}

var decompressCmd = &cobra.Command{
	Use:     fmt.Sprintf(CUT_OFF_COMMAND_PREFIX, "decompress"),
	Aliases: []string{fmt.Sprintf(CUT_OFF_COMMAND_PREFIX, "decompress")},
	Short:   "decompress files command, speed up decompress process by multiple parallelism",
	Long:    "decompress files command, speed up decompress process by multiple parallelism",
	Example: "co_decompress ./aa/ ./bb/ -parallelism 4",
	//Args:                       cobra.ExactArgs(),
	ArgAliases: nil,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := cmd.Flag("s").Value.String()
		targetPath := cmd.Flag("t").Value.String()
		parallelism := cmd.Flag("p").Value.String()
		var parallelismNumber int = 0
		var err error
		if len(parallelism) > 0 {
			parallelismNumber, err = strconv.Atoi(parallelism)
			if err != nil {
				panic(fmt.Sprintf("decode flag --p error: %v", err))
			}
		}
		handler, err := gzip.Get(sourcePath, targetPath, parallelismNumber, 0, false, true)
		if err != nil {
			log.Warn("init decompress env error: %v", err)
			return
		}
		log.Info("source path: %v, target path: %v, parallelism: %v", sourcePath, targetPath, handler.Parallelism)
		startTime := time.Now().Unix()
		err = handler.Decompress()
		if err != nil {
			log.Warn("decompress file error: %v", err)
			return
		}
		log.Info("decompress process finish, files: %v cost time: %vs", handler.FileNum, time.Now().Unix()-startTime)
	},
	RunE:                       nil,
	PostRun:                    nil,
	PostRunE:                   nil,
	PersistentPostRun:          nil,
	PersistentPostRunE:         nil,
	FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	CompletionOptions:          cobra.CompletionOptions{},
	TraverseChildren:           false,
	Hidden:                     false,
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
}
