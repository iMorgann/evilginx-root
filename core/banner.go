package core

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	VERSION = "2.4.2"
)

func putAsciiArt(s string) {
	for _, c := range s {
		d := string(c)
		switch string(c) {
		case "#":
			color.Set(color.BgRed)
			d = " "
		case "@":
			color.Set(color.BgBlack)
			d = " "
		case `:`:
			color.Set(color.BgGreen)
			d = " "
		case `$`:
			color.Set(color.BgYellow)
			d = " "
		case `/`:
			color.Set(color.BgBlue)
			d = " "
		case " ":
			color.Unset()
		}
		fmt.Print(d)
	}
	color.Unset()
}

func printLogo(s string) {
	for _, c := range s {
		d := string(c)
		switch string(c) {
		case "_":
			color.Set(color.FgWhite)
		case "\n":
			color.Unset()
		default:
			color.Set(color.FgHiBlack)
		}
		fmt.Print(d)
	}
	color.Unset()
}

func printUpdateName() {
	nameClr := color.New(color.FgHiRed)
	txt := nameClr.Sprintf("                 - --  Gone Phishing  -- -")
	fmt.Fprintf(color.Output, "%s", txt)
}

func printOneliner1() {
	handleClr := color.New(color.FgHiBlue)
	versionClr := color.New(color.FgGreen)
	textClr := color.New(color.FgHiBlack)
	spc := strings.Repeat(" ", 10-len(VERSION))
	txt := textClr.Sprintf("      by Kuba Gretzky (") + handleClr.Sprintf("@mrgretzky") + textClr.Sprintf(")") + textClr.Sprintf(", modified by (") + handleClr.Sprintf("@Anon_zz") + textClr.Sprintf(")") + spc + textClr.Sprintf("version ") + versionClr.Sprintf("%s", VERSION)
	fmt.Fprintf(color.Output, "%s", txt)
}

func printOneliner2() {
	textClr := color.New(color.FgHiBlack)
	red := color.New(color.FgRed)
	white := color.New(color.FgWhite)
	txt := red.Sprintf("                   Don't be evil") + white.Sprintf(" - ") + textClr.Sprintf("Anon_zz")
	fmt.Fprintf(color.Output, "%s", txt)
}

func Banner() {
	fmt.Println()

	putAsciiArt("             ##                                          ##           ")
	fmt.Println()
	putAsciiArt("             ####                                      ####           ")
	fmt.Println()
	putAsciiArt("             ######                                  ######            ")
	fmt.Println()
	putAsciiArt("              #######                              #######            ")
	fmt.Println()
	putAsciiArt("               ##########################################             ")
	fmt.Println()
	putAsciiArt("                 ######################################               ")
	fmt.Println()
	putAsciiArt("               ##########################################             ")
	fmt.Println()
	putAsciiArt("           ##################################################         ")
	fmt.Println()
	putAsciiArt("         #####################################################        ")
	fmt.Println()
	putAsciiArt(`        #########@@##################################@@########       `)
	fmt.Println()
	putAsciiArt("       ##########@@@@@@##########################@@@@@@#########      ")
	fmt.Println()
	putAsciiArt("      #############@@@@@@@#####################@@@@@@############     ")
	fmt.Println()
	putAsciiArt("     #############################################################    ")
	fmt.Println()
	putAsciiArt("    ###############################################################   ")
	fmt.Println()
	putAsciiArt("   #################################################################  ")
	fmt.Println()
	putAsciiArt("   #################################################################  ")
	fmt.Println()
	putAsciiArt("   $$$$$$$$$$$$$                      //////////////////////////////  ")
	fmt.Println()
	putAsciiArt("   $$$$$$$$$$$$$                      //////////////////////////////  ")
	fmt.Println()
	putAsciiArt("   $$$$$$$$$$$$$                      //////////////////////////////  ")
	fmt.Println()
	putAsciiArt("   $$$$$$$$$$$$$                      //////////////////////////////  ")
	fmt.Println()
	putAsciiArt("   $$$$$$$$$$$$:                                      //////////////  ")
	fmt.Println()
	putAsciiArt("    $$$$$$$$::::::                                   //////////////   ")
	fmt.Println()
	putAsciiArt("     $$$$:::::::::::                                //////////////    ")
	fmt.Println()
	putAsciiArt("      ::::::::::::::::                            ///////////////     ")
	fmt.Println()
	putAsciiArt("       :::::::::::::::::                        :://////////////      ")
	fmt.Println()
	putAsciiArt("         ::::::::::::::::::::::::::::::::::::::::::://////////        ")
	fmt.Println()
	putAsciiArt("           ::::::::::::::::::::::::::::::::::::::::::://////          ")
	fmt.Println()
	putAsciiArt("              :::::::::::::::::::::::::::::::::::::::::::             ")
	fmt.Println()
	putAsciiArt("                 :::::::::::::::::::::::::::::::::::::                ")
	fmt.Println()
	putAsciiArt("                    :::::::::::::::::::::::::::::::                   ")
	fmt.Println()
	putAsciiArt("                        :::::::::::::::::::::::                       ")
	fmt.Println()
	printLogo(`    ___________      __ __           __                   _             _                                 _`)
	fmt.Println()
	printLogo(`    \_   _____/__  _|__|  |    ____ |__| ____ ___  ___   | |__    ___  | |_  __ _  _   _   __ _  _ __  __| |`)
	fmt.Println()
	printLogo(`     |    __)_\  \/ /  |  |   / __ \|  |/    \\  \/  /   | '_ \  / _ \ | __|/ _  || | | | / _  || '__|/ _  |`)
	fmt.Println()
	printLogo(`     |        \\   /|  |  |__/ /_/  >  |   |  \>    <  _ | |_) || (_) || |_| (_| || |_| || (_| || |  | (_| |`)
	fmt.Println()
	printLogo(`    /_______  / \_/ |__|____/\___  /|__|___|  /__/\_ \(_)|_.__/  \___/  \__|\__, | \__,_| \__,_||_|   \__,_|`)
	fmt.Println()
	printLogo(`            \/              /_____/         \/      \/                      |___/ `)
	fmt.Println()
	printUpdateName()
	fmt.Println()
	printOneliner1()
	fmt.Println()
	printOneliner2()
	fmt.Println()
	fmt.Println()
}
