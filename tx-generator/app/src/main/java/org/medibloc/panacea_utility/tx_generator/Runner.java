package org.medibloc.panacea_utility.tx_generator;

import org.medibloc.panacea.*;
import org.medibloc.panacea.domain.TxResponse;
import org.medibloc.panacea.encoding.message.*;
import org.medibloc.panacea.encoding.message.did.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.security.NoSuchAlgorithmException;
import java.util.Collections;
import java.util.UUID;

public class Runner {
    private final static Logger LOG = LoggerFactory.getLogger(Runner.class);

    private final PanaceaApiRestClient client;
    private final Wallet wallet;

    private static final String DENOM = "umed";
    private static final String FEES = "1000000";
    private static final String GAS_LIMIT = "200000";
    private static final int SLEEP_MS = 5 * 1000;
    private static final String BROADCAST_MODE = "block";

    Runner(String lcdEndpoint, String mnemonic) {
        this.client = PanaceaApiClientFactory.newInstance().newRestClient(lcdEndpoint);
        this.wallet = Wallet.createWalletFromMnemonicCode(mnemonic, "panacea");
    }

    void run() throws PanaceaApiException, IOException, NoSuchAlgorithmException, TxGeneratorException, InterruptedException {
        String addr = this.wallet.getAddress();
        String topic = UUID.randomUUID().toString();
        createTopic(addr, topic);
        addWriter(addr, topic, addr);

        while (true) {
            sendToken(addr, "1", DENOM);
            Thread.sleep(SLEEP_MS);

            addRecord(addr, topic, "key-1".getBytes(), "value-1".getBytes(), addr);
            Thread.sleep(SLEEP_MS);

            createDid();
            Thread.sleep(SLEEP_MS);
        }
    }

    private void sendToken(String toAddr, String amount, String denom) throws PanaceaApiException, IOException, NoSuchAlgorithmException, TxGeneratorException {
        Coin coin = new Coin();
        coin.setAmount(amount);
        coin.setDenom(denom);

        MsgSend.Value value = new MsgSend.Value();
        value.setFromAddress(this.wallet.getAddress());
        value.setToAddress(toAddr);
        value.setAmount(Collections.singletonList(coin));

        MsgSend msg = new MsgSend();
        msg.setValue(value);

        executeTx(msg);
    }

    private void createTopic(String ownerAddr, String topic) throws PanaceaApiException, IOException, NoSuchAlgorithmException, TxGeneratorException {
        MsgCreateTopic.Value value = new MsgCreateTopic.Value();
        value.setTopicName(topic);
        value.setDescription("this is a topic description");
        value.setOwnerAddress(ownerAddr);

        MsgCreateTopic msg = new MsgCreateTopic();
        msg.setValue(value);

        executeTx(msg);
    }

    private void addWriter(String ownerAddr, String topic, String writerAddr) throws PanaceaApiException, IOException, NoSuchAlgorithmException, TxGeneratorException {
        MsgAddWriter.Value value = new MsgAddWriter.Value();
        value.setOwnerAddress(ownerAddr);
        value.setTopicName(topic);
        value.setMoniker("moniker");
        value.setDescription("this is a topic description");
        value.setWriterAddress(writerAddr);

        MsgAddWriter msg = new MsgAddWriter();
        msg.setValue(value);

        executeTx(msg);
    }

    private void addRecord(String ownerAddr, String topic, byte[] key, byte[] value, String writerAddr) throws PanaceaApiException, IOException, NoSuchAlgorithmException, TxGeneratorException {
        MsgAddRecord msg = new MsgAddRecord(ownerAddr, topic, key, value, writerAddr, this.wallet.getAddress());

        executeTx(msg);
    }

    private void createDid() throws NoSuchAlgorithmException, IOException, PanaceaApiException, TxGeneratorException {
        DidWallet didWallet = DidWallet.createRandomWallet();
        DidDocument didDocument = DidDocument.create(didWallet);

        MsgCreateDid msg = new MsgCreateDid(
                didDocument,
                didDocument.getVerificationMethods().get(0).getId(),
                didWallet,
                wallet.getAddress()
        );

        executeTx(msg);
    }

    private void executeTx(PanaceaTransactionMessage msg) throws PanaceaApiException, IOException, NoSuchAlgorithmException, TxGeneratorException {
        wallet.ensureWalletIsReady(client);

        StdTx tx = new StdTx(msg, new StdFee(DENOM, FEES, GAS_LIMIT), "");
        tx.sign(wallet);
        wallet.increaseAccountSequence();

        BroadcastReq req = new BroadcastReq(tx, BROADCAST_MODE);
        TxResponse res = client.broadcast(req);
        if (res.getCode() != 0) {
            throw new TxGeneratorException("TX failed: " + res);
        } else {
            LOG.info(res.toString());
        }
    }
}
