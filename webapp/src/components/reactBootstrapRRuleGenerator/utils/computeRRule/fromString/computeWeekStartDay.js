const computeWeekStartDay = (data, rruleObj) => {
    if (!rruleObj.wkst) {
        return data.options.weekStartsOnSunday;
    }
    return rruleObj.wkst === 6;
};

export default computeWeekStartDay;
