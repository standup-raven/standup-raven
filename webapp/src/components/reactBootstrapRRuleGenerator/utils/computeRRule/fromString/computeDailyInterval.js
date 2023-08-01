const computeDailyInterval = (data, rruleObj) => {
    if (rruleObj.freq !== 3) {
        return data.repeat.daily.interval;
    }

    return rruleObj.interval;
};

export default computeDailyInterval;
