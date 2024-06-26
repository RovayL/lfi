#include <stdbool.h>
#include <stdio.h>

#include "args.h"
#include "op.h"

bool amd64_parseinit();

extern struct op* ops;

typedef void (*PassFn)(struct op* op);

typedef struct pass {
    PassFn fn;
    bool disabled;
} Pass;

void amd64_specialpass(struct op*);
void amd64_branchpass(struct op*);
void amd64_loadspass(struct op*);
void amd64_storespass(struct op*);
void amd64_declpass(struct op*);
void amd64_syscallpass(struct op*);

static Pass passes[] = {
    (Pass) { .fn = &amd64_specialpass },
    (Pass) { .fn = &amd64_loadspass },
    (Pass) { .fn = &amd64_storespass },
    (Pass) { .fn = &amd64_branchpass },
    (Pass) { .fn = &amd64_declpass, .disabled = true },
    (Pass) { .fn = &amd64_syscallpass },
};

void amd64_display(FILE* output, struct op* ops);

static void
warnargs()
{
    if (args.poc)
        fprintf(stderr, "warning: --poc has no effect on amd64\n");
    if (args.noguardelim)
        fprintf(stderr, "warning: --no-guard-elim has no effect on amd64\n");
    if (args.sysexternal)
        fprintf(stderr, "warning: --sys-external has no effect on amd64\n");
    if (args.singlethread && args.boxtype != BOX_BUNDLEJUMPS)
        fprintf(stderr, "warning: --single-thread has no effect if --sandbox != bundle-jumps\n");
}

bool
amd64_rewrite(FILE* input, FILE* output)
{
    if (!amd64_parseinit()) {
        fprintf(stderr, "%s: parser failed to initialize\n", args.input);
        return false;
    }

    warnargs();

    const size_t npass = sizeof(passes) / sizeof(passes[0]);

    for (size_t i = 0; i < npass; i++) {
        if (args.boxtype < BOX_FULL && passes[i].fn == &amd64_loadspass)
            passes[i].disabled = true;
        if (args.boxtype < BOX_STORES && passes[i].fn == &amd64_storespass)
            passes[i].disabled = true;
        if (args.boxtype < BOX_STORES && passes[i].fn == &amd64_specialpass)
            passes[i].disabled = true;
        if (args.boxtype < BOX_BUNDLEJUMPS && passes[i].fn == &amd64_branchpass)
            passes[i].disabled = true;
        if (args.decl && passes[i].fn == &amd64_declpass)
            passes[i].disabled = false;
    }

    for (size_t i = 0; i < npass; i++) {
        if (passes[i].disabled)
            continue;
        struct op* op = ops;
        while (op) {
            struct op* next = op->next;
            passes[i].fn(op);
            op = next;
        }
    }

    amd64_display(output, ops);

    return true;
}
