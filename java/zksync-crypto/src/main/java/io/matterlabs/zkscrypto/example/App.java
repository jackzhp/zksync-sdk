/*
 * This Java source file was generated by the Gradle 'init' task.
 */
package io.matterlabs.zkscrypto.example;

import io.matterlabs.zkscrypto.lib.ZksCrypto;
import io.matterlabs.zkscrypto.lib.entity.ZksPackedPublicKey;
import io.matterlabs.zkscrypto.lib.entity.ZksPrivateKey;
import io.matterlabs.zkscrypto.lib.entity.ZksPubkeyHash;
import io.matterlabs.zkscrypto.lib.entity.ZksSignature;
import io.matterlabs.zkscrypto.lib.exception.ZksMusigTooLongException;
import io.matterlabs.zkscrypto.lib.exception.ZksSeedTooShortException;

import org.web3j.utils.Numeric;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;

public class App {

    public static void main(String[] args) {
        ZksCrypto crypto = ZksCrypto.load();

        byte[] seed = Arrays.copyOf(new byte[0], 32);
        byte[] msg = "hello".getBytes(StandardCharsets.UTF_8);

        try {
            ZksPrivateKey privateKey = crypto.generatePrivateKey(seed);
            ZksPackedPublicKey publicKey = crypto.getPublicKey(privateKey);
            ZksPubkeyHash pubkeyHash = crypto.getPublicKeyHash(publicKey);
            ZksSignature signature = crypto.signMessage(privateKey, msg);
            System.out.printf("Seed: %s\n", Numeric.toHexString(seed));
            System.out.printf("Private key: %s\n", Numeric.toHexString(privateKey.getData()));
            System.out.printf("Public key: %s\n", Numeric.toHexString(publicKey.getData()));
            System.out.printf("Public key hash: %s\n", Numeric.toHexString(pubkeyHash.getData()));
            System.out.printf("Signature: %s\n", Numeric.toHexString(signature.getData()));
        } catch (ZksSeedTooShortException | ZksMusigTooLongException e) {
            System.err.println(e);
        }
    }

}
