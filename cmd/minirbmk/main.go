// main.go - Main for the minirbmk example
// SPDX-License-Identifier: GPL-3.0-or-later

// The minirbmk command shows how to write a [clip] based
// command line tool using nested subcommands.
package main

import "github.com/bassosimone/clip"

// configurable for testing
var (
	// env is the execution environment to use
	env = clip.NewStdlibExecEnv()
)

func main() {
	// Define the overall suite version
	const version = "0.1.0"

	// Create the curl leaf command
	curlCmd := &clip.LeafCommand[*clip.StdlibExecEnv]{
		BriefDescriptionText: "Utility to transfer URLs.",
		RunFunc:              curlMain,
	}

	// Create the dig leaf command
	digCmd := &clip.LeafCommand[*clip.StdlibExecEnv]{
		BriefDescriptionText: "Utility to query the DNS.",
		RunFunc:              digMain,
	}

	// Create the 'git clone' leaf command
	gitCloneCmd := &clip.LeafCommand[*clip.StdlibExecEnv]{
		BriefDescriptionText: "Clone a repository.",
		RunFunc:              gitCloneMain,
	}

	// Create the 'git init' leaf command.
	gitInitCmd := &clip.LeafCommand[*clip.StdlibExecEnv]{
		BriefDescriptionText: "Init a repository.",
		RunFunc:              gitInitMain,
	}

	// Create the git subcommand
	gitCmd := &clip.DispatcherCommand[*clip.StdlibExecEnv]{
		BriefDescriptionText: "Utility to manage repositories.",
		Commands: map[string]clip.Command[*clip.StdlibExecEnv]{
			"clone": gitCloneCmd,
			"init":  gitInitCmd,
		},
		ErrorHandling: clip.ExitOnError,
		Version:       version,
	}

	// Create the root command
	rootCmd := &clip.RootCommand[*clip.StdlibExecEnv]{
		// Use a dispatcher dispatching to `git`, `curl`, and `dig`.
		Command: &clip.DispatcherCommand[*clip.StdlibExecEnv]{

			// This text is printed when help is requested
			BriefDescriptionText: "A collection of UNIX command line tools.",

			// Configure the dispatcher to dispatch by name
			Commands: map[string]clip.Command[*clip.StdlibExecEnv]{
				"curl": curlCmd,
				"dig":  digCmd,
				"git":  gitCmd,
			},

			// Cause the dispatcher to call [os.Exit] on error
			ErrorHandling: clip.ExitOnError,

			// Automatically define --version and the version subcommand
			Version: version,
		},

		// Automatic signals handling: SIGINT and SIGTERM will
		// cancel the context passed to leaf commands.
		AutoCancel: true,
	}

	// Execute the root command
	rootCmd.Main(env)
}
