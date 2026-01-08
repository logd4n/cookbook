/*  # # #   # # #  # # #  # # # #       # # # #     #       #        # # #   # # # #
  #     # #     # #    # #             #    # #   #        #       # * * #  #
 #       #     # #    # # # # #       # # #   # #         #       #  *  #  #  # # #
#     # #     # #    # #             #    #   #          #       # * * #  #      #
# # #   # # #  # # #  # # # #       # # #    #          # # # #  # # #    # # # # */

//ColorsByLOG.go --version 1.1.0

package colors

const (
	DefaultColor string = "\x1b[0m"

	Text_Black  string = "\x1b[30m"
	Text_Red    string = "\x1b[31m"
	Text_Green  string = "\x1b[32m"
	Text_Yellow string = "\x1b[33m"
	Text_Blue   string = "\x1b[34m"
	Text_Purple string = "\x1b[35m"
	Text_Cyan   string = "\x1b[36m"
	Text_White  string = "\x1b[37m"

	Attribute_Bold       string = "\x1b[1m"
	Attribute_Italic     string = "\x1b[3m"
	Attribute_Underlined string = "\x1b[4m"
	Attribute_Invisible  string = "\x1b[8m"

	Backgrond_Black  string = "\x1b[40m"
	Backgrond_Red    string = "\x1b[41m"
	Backgrond_Green  string = "\x1b[42m"
	Backgrond_Yellow string = "\x1b[43m"
	Backgrond_Blue   string = "\x1b[44m"
	Backgrond_Purple string = "\x1b[45m"
	Backgrond_Cyan   string = "\x1b[46m"
	Backgrond_White  string = "\x1b[47m"
)

func SetColor(_colorName string) {
	print(_colorName)
}

func ResetColor() {
	print(DefaultColor)
}
