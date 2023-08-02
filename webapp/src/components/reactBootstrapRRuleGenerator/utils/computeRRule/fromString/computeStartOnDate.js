const computeStartOnDate = (data, rruleObj) => {
    if (!rruleObj.dtstart) {
        return data.start.onDate.date;
    }

    return rruleObj.dtstart;
};
export default computeStartOnDate;
