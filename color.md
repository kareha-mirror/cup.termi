# Color Codes in termi

[Japanese version](color.ja.md)

Terminal emulators generally support three color modes: **16-color**, **256-color**, and **TrueColor**.
Each mode uses a different method of specifying colors.

The **termi** terminal library hides these differences.
Colors can be specified using any supported notation, and termi automatically maps them to the closest available color for the current terminal.

Text displayed in a terminal consists of a **foreground** color and an optional **background** color.
In termi, these are specified together by separating them with a comma `,`.

## 16-Color Mode

In 16-color mode, colors can be specified either by their numeric value (`0`-`15`) or by a **color name**.

The predefined colors are:

* 0: black
* 1: red
* 2: green
* 3: yellow
* 4: blue
* 5: magenta
* 6: cyan
* 7: white
* 8: bright black
* 9: bright red
* 10: bright green
* 11: bright yellow
* 12: bright blue
* 13: bright magenta
* 14: bright cyan
* 15: bright white

Special name `default` refers to default color of the terminal.

Color names are case-insensitive, and word separators are optional. For example, all of the following refer to color `10` (**bright green**):

* bright green
* bright-green
* brightgreen
* Bright Green
* BrightGreen

## 256-Color Mode

In 256-color mode, colors are specified by a number from `0` to `255`.

Values `0`-`15` correspond to the standard 16-color palette.

You can display the complete 256-color palette with:

```text
termi color
```

Values `16`-`231` form the color cube, while `232`-`255` are grayscale shades.

Refer to the color chart for the complete palette.

## TrueColor Mode

In TrueColor mode, colors are specified as a six-digit hexadecimal **RGB** value.

For example, a color with red `ff`, green `80`, and blue `40` is written as:

```text
ff8040
```

## Specifying Foreground and Background Colors

Foreground and background colors are specified together, separated by a comma.

For example, foreground color `252` with background color `235` is written as:

```text
252,235
```
