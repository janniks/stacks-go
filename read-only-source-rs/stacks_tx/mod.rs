use neon::prelude::*;
use sha2::{Digest, Sha512_256};
use std::io::Cursor;

use crate::hex::encode_hex;
use crate::neon_util::*;

use self::deserialize::StacksTransaction;
mod deserialize;
mod neon_encoder;

pub fn decode_transaction(mut cx: FunctionContext) -> JsResult<JsObject> {
    let (tx, tx_id_bytes) = arg_as_bytes(&mut cx, 0, |val_bytes| {
        let mut cursor = Cursor::new(val_bytes);
        let tx = StacksTransaction::deserialize(&mut cursor)
            .or_else(|e| Err(format!("Failed to decode transaction: {:?}\n", &e)))?;
        let tx_id_bytes = Sha512_256::digest(val_bytes);
        Ok((tx, tx_id_bytes))
    })
    .or_else(|e| cx.throw_error(e))?;

    let tx_json_obj = cx.empty_object();

    let tx_id = cx.string(encode_hex(&tx_id_bytes));
    tx_json_obj.set(&mut cx, "tx_id", tx_id)?;

    tx.neon_js_serialize(&mut cx, &tx_json_obj, &())?;
    Ok(tx_json_obj)
}
