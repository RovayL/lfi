#pragma once

#include <stdbool.h>

struct op {
    char* text;

    bool insn;
    bool relocated;
    char* label;

    bool shortbr;
    char* replace;
    char* target;

    struct op* next;
    struct op* prev;
};

struct op;

struct op* mktbz(char* tbz, char* reg, char* imm, char* label);

struct op* mklabel(char* name);

struct op* mkinsn(char* fmt, ...);

struct op* mkdirective(char* text);

void opremove(struct op* op);

void opfree(struct op* op);

void oplocate(struct op* op);

struct op* opinsert(struct op* op);
