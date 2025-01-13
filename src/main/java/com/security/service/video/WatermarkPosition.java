package com.security.service.video;

public enum WatermarkPosition {
    TOP_LEFT("10:10"),
    TOP_RIGHT("main_w-overlay_w-10:10"),
    BOTTOM_LEFT("10:main_h-overlay_h-10"),
    BOTTOM_RIGHT("main_w-overlay_w-10:main_h-overlay_h-10"),
    CENTER("(main_w-overlay_w)/2:(main_h-overlay_h)/2");
    
    private final String coordinates;
    
    WatermarkPosition(String coordinates) {
        this.coordinates = coordinates;
    }
    
    public String getCoordinates() {
        return coordinates;
    }
} 