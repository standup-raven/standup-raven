const computeHourlyInterval = (data, rruleObj) => {
    if (rruleObj.freq !== 4) {
        return data.repeat.daily.interval;
    }

    return rruleObj.interval;
};

export default computeHourlyInterval;
