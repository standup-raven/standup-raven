const computeYearlyOnMonthday = (data, rruleObj) => {
    if (rruleObj.freq !== 0 || !rruleObj.bymonthday) {
        return data.repeat.yearly.on.day;
    }

    if (typeof rruleObj.bymonthday === 'number') {
        return rruleObj.bymonthday;
    }

    return rruleObj.bymonthday[0];
};

export default computeYearlyOnMonthday;
