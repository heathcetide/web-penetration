package com.security.config;

public class VideoConfig {
    private boolean drmEnabled = false;
    private String drmLicenseServer;
    private String widevineKeyServer;
    private String playreadyKeyServer;
    private int maxRetries = 3;
    private long retryDelay = 1000;
    private String tempDir = "temp/video/";
    private EncryptionConfig encryption = new EncryptionConfig();

    public static class EncryptionConfig {
        private boolean enabled = true;
        private String algorithm = "AES";
        private String keySize = "128";
        private String mode = "CBC";
        private String padding = "PKCS5Padding";

        public boolean isEnabled() {
            return enabled;
        }

        public void setEnabled(boolean enabled) {
            this.enabled = enabled;
        }

        public String getAlgorithm() {
            return algorithm;
        }

        public void setAlgorithm(String algorithm) {
            this.algorithm = algorithm;
        }

        public String getKeySize() {
            return keySize;
        }

        public void setKeySize(String keySize) {
            this.keySize = keySize;
        }

        public String getMode() {
            return mode;
        }

        public void setMode(String mode) {
            this.mode = mode;
        }

        public String getPadding() {
            return padding;
        }

        public void setPadding(String padding) {
            this.padding = padding;
        }
    }

    public boolean isDrmEnabled() {
        return drmEnabled;
    }

    public void setDrmEnabled(boolean drmEnabled) {
        this.drmEnabled = drmEnabled;
    }

    public String getDrmLicenseServer() {
        return drmLicenseServer;
    }

    public void setDrmLicenseServer(String drmLicenseServer) {
        this.drmLicenseServer = drmLicenseServer;
    }

    public String getWidevineKeyServer() {
        return widevineKeyServer;
    }

    public void setWidevineKeyServer(String widevineKeyServer) {
        this.widevineKeyServer = widevineKeyServer;
    }

    public String getPlayreadyKeyServer() {
        return playreadyKeyServer;
    }

    public void setPlayreadyKeyServer(String playreadyKeyServer) {
        this.playreadyKeyServer = playreadyKeyServer;
    }

    public int getMaxRetries() {
        return maxRetries;
    }

    public void setMaxRetries(int maxRetries) {
        this.maxRetries = maxRetries;
    }

    public long getRetryDelay() {
        return retryDelay;
    }

    public void setRetryDelay(long retryDelay) {
        this.retryDelay = retryDelay;
    }

    public String getTempDir() {
        return tempDir;
    }

    public void setTempDir(String tempDir) {
        this.tempDir = tempDir;
    }

    public EncryptionConfig getEncryption() {
        return encryption;
    }

    public void setEncryption(EncryptionConfig encryption) {
        this.encryption = encryption;
    }
}