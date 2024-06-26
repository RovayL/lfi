
extern crate alloc;

use alloc::format;
use verifier::Verifier;

mod inst;
mod verifier;
mod autogen;
mod reg;

#[no_mangle]
pub extern "C" fn lfi_verify_insn(insn: u32) -> bool {
    let mut verif = Verifier {
        failed: false,
        message: None,
    };

    if let Ok(decoded) = bad64::decode(insn, 0) {
        verif.check_insn(&decoded);
        return !verif.failed;
    }
    return false;
}

#[no_mangle]
pub extern "C" fn lfi_verify_bytes(raw_bytes: *const u8, size: usize, error: *const ()) -> bool {
    let bytes = unsafe { core::slice::from_raw_parts(raw_bytes, size) };

    let fnptr: Option<fn(*const u8, usize)> = if error != core::ptr::null() {
        Some(unsafe { core::mem::transmute(error) })
    } else {
        None
    };

    let mut verif = Verifier {
        failed: false,
        message: fnptr,
    };

    let mut iter = bad64::disasm(bytes, 0);
    while let Some(maybe_decoded) = iter.next() {
        match maybe_decoded {
            Ok(inst) => {
                verif.check_insn(&inst);
            }
            Err(e) => {
                if let Some(message) = fnptr {
                    let s = format!("{:x}: unknown instruction: {}", e.address(), e);
                    message(s.as_ptr(), s.len());
                } else {
                    // if no messaging, return immediately
                    return false;
                }
                verif.failed = true;
            }
        }
    }

    return !verif.failed;
}
