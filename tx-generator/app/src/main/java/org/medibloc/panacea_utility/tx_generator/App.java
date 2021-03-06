/*
 * This Java source file was generated by the Gradle 'init' task.
 */
package org.medibloc.panacea_utility.tx_generator;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class App {
    private final static Logger LOG = LoggerFactory.getLogger(App.class);

    public static void main(String[] args) {
        String mnemonic = System.getenv("MNEMONIC");
        if (mnemonic == null) {
            LOG.error("MNEMONIC env var must be specified");
            System.exit(1);
        }

        String lcdEndpoint = System.getenv("LCD_ENDPOINT");
        if (lcdEndpoint == null) {
            LOG.error("LCD_ENDPOINT env var must be specified");
            System.exit(1);
        }
        LOG.info("LCD_ENDPOINT: " + lcdEndpoint);

        Runner runner = new Runner(lcdEndpoint, mnemonic);
        try {
            runner.run();
        } catch (Exception e) {
            LOG.error(e.getLocalizedMessage());
            System.exit(1);
        }
    }
}
