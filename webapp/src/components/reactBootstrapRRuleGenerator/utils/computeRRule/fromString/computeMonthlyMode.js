const computeMonthlyMode = (data, rruleObj) => {
    if (rruleObj.freq !== 1) {
        return data.repeat.monthly.mode;
    }

    if (rruleObj.bymonthday) {
        return 'on';
    }

    return 'on the';
};

export default computeMonthlyMode;
