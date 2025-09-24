package main

import (
    "flag"
    "fmt"
    "os"
)

var appVersion = "0.5.0" // Default version, can be overridden via ldflags

var (
    debugFlag   bool
    verboseFlag bool
    noColorFlag bool
    helpFlag    bool
    versionFlag bool
    configFile  string
    gitFlag     bool
)

// Color codes for terminal output
type Colors struct {
    Reset  string
    Red    string
    Green  string
    Yellow string
    Blue   string
    Purple string
    Cyan   string
    Bold   string
}

var colors Colors

func init() {
    flag.BoolVar(&debugFlag, "debug", false, "enable debug output")
    flag.BoolVar(&debugFlag, "d", false, "enable debug output (shorthand)")
    flag.BoolVar(&verboseFlag, "verbose", false, "verbose output")
    flag.BoolVar(&verboseFlag, "v", false, "verbose output (shorthand)")
    flag.BoolVar(&noColorFlag, "no-color", false, "disable colored output")
    flag.BoolVar(&helpFlag, "help", false, "print help and exit")
    flag.BoolVar(&helpFlag, "h", false, "print help and exit (shorthand)")
    flag.BoolVar(&versionFlag, "version", false, "print version and exit")
    flag.BoolVar(&versionFlag, "V", false, "print version and exit (shorthand)")
    flag.StringVar(&configFile, "config", "", "specify custom .project.yml configuration file path")
    flag.BoolVar(&gitFlag, "git", false, "force use of git-based detection (ignore .project.yml)")
}

func setupColors() {
    if noColorFlag || !isTerminal() {
        colors = Colors{
            Reset:  "",
            Red:    "",
            Green:  "",
            Yellow: "",
            Blue:   "",
            Purple: "",
            Cyan:   "",
            Bold:   "",
        }
    } else {
        colors = Colors{
            Reset:  "\033[0m",
            Red:    "\033[31m",
            Green:  "\033[32m",
            Yellow: "\033[33m",
            Blue:   "\033[34m",
            Purple: "\033[35m",
            Cyan:   "\033[36m",
            Bold:   "\033[1m",
        }
    }
}

func isTerminal() bool {
    // Check if stderr is a terminal
    fileInfo, err := os.Stderr.Stat()
    if err != nil {
        return false
    }
    return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func printError(format string, args ...interface{}) {
    message := fmt.Sprintf(format, args...)
    fmt.Fprintf(os.Stderr, "%s%sERROR%s%s: %s%s\n", 
        colors.Red, colors.Bold, colors.Reset, colors.Red, message, colors.Reset)
}

func printWarning(format string, args ...interface{}) {
    message := fmt.Sprintf(format, args...)
    fmt.Fprintf(os.Stderr, "%sWARNING: %s%s\n", colors.Purple, message, colors.Reset)
}

func printDebug(format string, args ...interface{}) {
    if debugFlag {
        message := fmt.Sprintf(format, args...)
        fmt.Fprintf(os.Stderr, "%s%s#DEBUG%s%s: %s%s\n", 
            colors.Yellow, colors.Bold, colors.Reset, colors.Yellow, message, colors.Reset)
    }
}

func printInfo(format string, args ...interface{}) {
    if verboseFlag {
        message := fmt.Sprintf(format, args...)
        fmt.Fprintf(os.Stderr, "%sINFO: %s%s\n", colors.Green, message, colors.Reset)
    }
}

func printSuccess(format string, args ...interface{}) {
    message := fmt.Sprintf(format, args...)
    fmt.Fprintf(os.Stderr, "%s%sSUCCESS%s%s: %s%s\n", 
        colors.Green, colors.Bold, colors.Reset, colors.Green, message, colors.Reset)
}

func printHelp() {
    fmt.Printf(`version %s - semantic version parsing, validation, and ordering utility

Usage:
    version [options] command [arguments]

Options:
    -h, --help        print this help and exit
    -V, --version     print version and exit
    -v, --verbose     verbose output
    -d, --debug       debug output
    --no-color        disable colored output
    --config FILE     specify custom .project.yml configuration file path
    --git             force use of git-based detection (ignore .project.yml)

Commands:
    project           print project name from git remote
    module            print module name from git remote
    modules           print all module names from .project.yml or single git module name
    version           print project version from git tags
    release           print project release number
    full              print full project name-version-release
    check [version]   validate version string (uses current git version if not specified)
    check-greatest [version] check if version is greatest among all tags
    type [version]    print version type (release, prerelease, postrelease, intermediate)
    build-type [version] print CMake build type (Release/Debug) based on version type
    bump [version] [type] bump version with specified type (smart, major, minor, patch, pre, alpha, beta, rc, fix, next, post, feat)
    sort              sort version strings from stdin
    platform          print current platform (GOOS value)
    arch              print current architecture (GOARCH value)
    os                print current operating system (user-friendly format)
    os_version        print current operating system version
    cpu               print number of logical CPUs

Examples:
    version check 1.2.3
    version check-greatest
    version type 1.2.3-alpha.1
    version bump 1.2.3 major
    version bump 1.2.3 alpha
    version platform
    version arch
    version os
    version os_version
    version cpu
    echo "1.2.3 1.2.4 1.2.3-alpha" | version sort
`, appVersion)
}

func main() {
    flag.Parse()
    setupColors()

    if helpFlag {
        printHelp()
        os.Exit(0)
    }

    if versionFlag {
        fmt.Println(appVersion)
        os.Exit(0)
    }

    // Validate conflicting flags
    if configFile != "" && gitFlag {
        printError("cannot use both --config and --git flags simultaneously")
        os.Exit(1)
    }

    args := flag.Args()
    if len(args) == 0 {
        printError("no command specified")
        printHelp()
        os.Exit(1)
    }

    command := args[0]
    commandArgs := args[1:]

    printDebug("Executing command: %s with args: %v", command, commandArgs)

    var err error
    var result string

    switch command {
    case "project":
        result, err = getProject()
    case "module":
        result, err = getModule()
    case "modules":
        result, err = getModules()
    case "version":
        result, err = getVersion()
    case "release":
        result, err = getRelease()
    case "full":
        result, err = getFull()
    case "check":
        if len(commandArgs) > 0 {
            err = checkVersion(commandArgs[0])
        } else {
            version, e := getVersion()
            if e != nil {
                err = e
            } else {
                err = checkVersion(version)
            }
        }
    case "check-greatest":
        if len(commandArgs) > 0 {
            result, err = checkGreatest(commandArgs[0])
        } else {
            version, e := getVersion()
            if e != nil {
                err = e
            } else {
                result, err = checkGreatest(version)
            }
        }
    case "type":
        if len(commandArgs) > 0 {
            result, err = getVersionType(commandArgs[0])
        } else {
            version, e := getVersion()
            if e != nil {
                err = e
            } else {
                result, err = getVersionType(version)
            }
        }
    case "build-type":
        if len(commandArgs) > 0 {
            result, err = getBuildType(commandArgs[0])
        } else {
            version, e := getVersion()
            if e != nil {
                err = e
            } else {
                result, err = getBuildType(version)
            }
        }
    case "sort":
        result, err = sortVersions()
    case "bump":
        if len(commandArgs) > 0 && (commandArgs[0] == "help" || commandArgs[0] == "--help" || commandArgs[0] == "-h") {
            printBumpHelp()
            os.Exit(0)
        }
        
        // Validate arguments
        if err := validateBumpArgs(commandArgs); err != nil {
            printError("%v", err)
            printBumpHelp()
            os.Exit(1)
        }
        
        // Get version to bump
        versionToBump, err := getBumpVersion(commandArgs)
        if err != nil {
            printError("%v", err)
            os.Exit(1)
        }
        
        // Get bump type
        bumpType := getBumpType(commandArgs)
        
        // Perform bump
        result, err = bumpVersion(versionToBump, bumpType)
    case "platform":
        result, err = getPlatform()
    case "arch":
        result, err = getArch()
    case "os":
        result, err = getOS()
    case "os_version":
        result, err = getOSVersion()
    case "cpu":
        result, err = getCPU()
    default:
        printError("unknown command: %s", command)
        printHelp()
        os.Exit(1)
    }

    if err != nil {
        printError("%v", err)
        os.Exit(1)
    }

    if result != "" {
        fmt.Println(result)
    }
}