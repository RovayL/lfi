#pragma once

#include <stdbool.h>

enum cfi {
    CFI_BUNDLE16,
    CFI_BUNDLE32,
    CFI_HW,
};

enum boxtype {
    BOX_NONE        = 0,
    BOX_BUNDLEJUMPS = 1,
    BOX_JUMPS       = 2,
    BOX_STORES      = 3,
    BOX_FULL        = 4,
};

enum meter {
    METER_NONE,
    METER_BRANCH,
    METER_FP,
};

enum flags {
    FLAGS_NONE,
    FLAGS_CLANG,
    FLAGS_GCC,
};

struct arguments {
    char* input;
    char* output;
    bool poc;
    bool sysexternal;
    bool noguardelim;
    bool nosegue;
    char* arch;
    enum cfi cfi;
    enum boxtype boxtype;
    bool singlethread;
    bool decl;
    enum meter meter;
    enum flags flags;
};

extern struct arguments args;
