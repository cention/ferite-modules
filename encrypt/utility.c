#include "encrypt_header.h"

MCRYPT Ferite_Encrypt_Open_MCrypt(void *key, void *IV) {
	MCRYPT td = mcrypt_module_open(MCRYPT_BLOWFISH, NULL, MCRYPT_CFB, NULL);
	mcrypt_generic_init(td, key, strlen(key), IV);
	return td;
}

void Ferite_Encrypt_Close_MCrypt(MCRYPT td) {
	mcrypt_generic_deinit(td);
	mcrypt_module_close(td);
}
