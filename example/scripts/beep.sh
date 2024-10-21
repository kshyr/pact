#!/bin/bash

# WAV file header and data binary strings
HEADER_START="\\x52\\x49\\x46\\x46\\x24\\x00\\x08\\x00\\x57\\x41\\x56\\x45\\x66\\x6d\\x74\\x20\\x10\\x00\\x00\\x00\\x01\\x00\\x01\\x00"
RATE="\\x40\\x1f\\x00\\x00\\x40\\x1f\\x00\\x00"
HEADER_END="\\x01\\x00\\x08\\x00\\x64\\x61\\x74\\x61\\x00\\x00\\x08\\x00"
DATA="\\x80\\x26\\x00\\x26\\x7F\\xD9\\xFF\\xD9"

generate_beep() {
    printf "%b" "$HEADER_START" > /tmp/sinewave.wav
    printf "%b" "$RATE" >> /tmp/sinewave.wav
    printf "%b" "$HEADER_END" >> /tmp/sinewave.wav
    for _ in {0..5}; do
        DATA=$DATA$DATA
    done
    printf "%b" "$DATA" >> /tmp/sinewave.wav
}

FREQ=$1
BEEP_DURATION_MS=$2

if [[ -z "$FREQ" ]] || [[ "$FREQ" -gt 1047 ]] || [[ "$FREQ" -lt 110 ]]; then
    echo "Frequency must be between 131Hz (C3) and 988Hz (C6)."
    exit 1
fi

# convert frequency to the appropriate sample rate
SAMPLE=$(( 8 * FREQ ))
SAMPLE=$(printf "%04x" "$SAMPLE")
RATE="\\x${SAMPLE:2:2}\\x${SAMPLE:0:2}\\x00\\x00\\x${SAMPLE:2:2}\\x${SAMPLE:0:2}\\x00\\x00"

generate_beep

# afplay (Mac) or aplay (Linux)
if [[ $(uname) == "Darwin" ]]; then
    afplay -t $((BEEP_DURATION_MS / 1000)) /tmp/sinewave.wav 2>/dev/null
else
    aplay -d $((BEEP_DURATION_MS / 1000)) /tmp/sinewave.wav
fi
