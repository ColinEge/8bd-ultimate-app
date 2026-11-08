export interface Key {
    id: string;
    colSpan?: number;
    rowSpan?: number;
}

/* -------------------------
     LEFT TOP (Esc + F row)
     ------------------------- */

export const fRow: Key[] = [
    { id: "0x29" }, // Esc
    { id: "0x3a" }, // F1
    { id: "0x3b" }, // F2
    { id: "0x3c" }, // F3
    { id: "0x3d" }, // F4
    { id: "0x3e" }, // F5
    { id: "0x3f" }, // F6
    { id: "0x40" }, // F7
    { id: "0x41" }, // F8
    { id: "0x42" }, // F9
    { id: "0x43" }, // F10
    { id: "0x44" }, // F11
    { id: "0x45" }, // F12
];

/* -------------------------
   LEFT MAIN (Typing cluster)
   15 grid columns
   ------------------------- */

export const typingRows: Key[][] = [
    [
        { id: "0x35" }, // `
        { id: "0x1e" }, // 1
        { id: "0x1f" }, // 2
        { id: "0x20" }, // 3
        { id: "0x21" }, // 4
        { id: "0x22" }, // 5
        { id: "0x23" }, // 6
        { id: "0x24" }, // 7
        { id: "0x25" }, // 8
        { id: "0x26" }, // 9
        { id: "0x27" }, // 0
        { id: "0x2d" }, // -
        { id: "0x2e" }, // =
        { id: "0x2a", colSpan: 2 }, // Backspace
    ],
    [
        { id: "0x2b", colSpan: 1.25 }, // Tab
        { id: "0x14" }, // Q
        { id: "0x1a" }, // W
        { id: "0x08" }, // E
        { id: "0x15" }, // R
        { id: "0x17" }, // T
        { id: "0x1c" }, // Y
        { id: "0x18" }, // U
        { id: "0x0c" }, // I
        { id: "0x12" }, // O
        { id: "0x13" }, // P
        { id: "0x2f" }, // [
        { id: "0x30" }, // ]
        { id: "0x31", colSpan: 1.75 }, // \
    ],
    [
        { id: "0x39", colSpan: 1.5 }, // Capslock
        { id: "0x04" }, // A
        { id: "0x16" }, // S
        { id: "0x07" }, // D
        { id: "0x09" }, // F
        { id: "0x0a" }, // G
        { id: "0x0b" }, // H
        { id: "0x0d" }, // J
        { id: "0x0e" }, // K
        { id: "0x0f" }, // L
        { id: "0x33" }, // ;
        { id: "0x34" }, // '
        { id: "0x28", colSpan: 2.5 }, // Enter
    ],
    [
        { id: "0x65", colSpan: 2.25 }, // leftshift
        { id: "0x1d" }, // Z
        { id: "0x1b" }, // X
        { id: "0x06" }, // C
        { id: "0x19" }, // V
        { id: "0x05" }, // B
        { id: "0x11" }, // N
        { id: "0x10" }, // M
        { id: "0x36" }, // ,
        { id: "0x37" }, // .
        { id: "0x38" }, // /
        { id: "0x69", colSpan: 2.75 },  // rightshift
    ],
    [
        { id: "0x64", colSpan: 1.25 }, // Left Ctrl
        { id: "0x67", colSpan: 1.25 }, // Win
        { id: "0x66", colSpan: 1.25 }, // Left Alt
        { id: "0x2c", colSpan: 6.5 }, // Space
        { id: "0x6a", colSpan: 1.25 }, // Alt
        { id: "0xe8", colSpan: 1.25 }, // Super B
        { id: "0xe9", colSpan: 1.25 }, // Super A
        { id: "0x68" }, // Right Ctrl
    ],
];

/* -------------------------
   MIDDLE (Nav cluster)
   ------------------------- */

export const navTop: Key[][] = [
    [
        { id: "0x46" }, // prtsc
        { id: "0x47" }, // scrlk
        { id: "0x48" }, // pause
    ],
];

export const navMain: Key[][] = [
    [
        { id: "0x49" }, // insert
        { id: "0x4a" }, // home
        { id: "0x4b" }, // pageup
    ],
    [
        { id: "0x4c" }, // delete
        { id: "0x4d" }, // end
        { id: "0x4e" }, // pagedown
    ],
];

// Keys with a 0x00 ID should take up space but not be visible or interactive
export const arrows: Key[][] = [
    [
        { id: "" }, 
        { id: "0x52" }, // up
        { id: "" },
    ],
    [
        { id: "0x50" }, // left
        { id: "0x51" }, // down
        { id: "0x4f" }, // right
    ],
];

/* -------------------------
   NUMPAD (Right side)
   ------------------------- */

// Keys with a 0x00 ID should take up space but not be visible or interactive
export const numpadTop: Key[][] = [
    [
        { id: "" },
        { id: "" },
        { id: "" },
        { id: "" },
    ],
];

export const numpad: Key[] = [
    { id: "0x53" }, // numlock
    { id: "0x54" }, // kpslash
    { id: "0x55" }, // kpasterisk
    { id: "0x56" }, // kpminus

    { id: "0x5f" }, // kp7
    { id: "0x60" }, // kp8
    { id: "0x61" }, // kp9
    { id: "0x57", rowSpan: 2 }, // kpplus

    { id: "0x5c" }, // kp4
    { id: "0x5d" }, // kp5
    { id: "0x5e" }, // kp6

    { id: "0x59" }, // kp1
    { id: "0x5a" }, // kp2
    { id: "0x5b" }, // kp3
    { id: "0x58", rowSpan: 2 }, // kpenter

    { id: "0x62", colSpan: 2 }, // kp0
    { id: "0x63" }, // kpdot
];