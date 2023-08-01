const computeEndAfter = (data, rruleObj) => {
    if (!rruleObj.count && rruleObj.count !== 0) {
        return data.end.after;
    }

    return rruleObj.count;
};

export default computeEndAfter;
