$(function () {
    $.fn.inputFilter = function (inputFilter) {
        return this.on("input", function () {
            if (inputFilter(this.value)) {
                this.oldValue = this.value;
                this.oldSelectionStart = this.selectionStart;
                this.oldSelectionEnd = this.selectionEnd;
            } else if (this.hasOwnProperty("oldValue")) {
                this.value = this.oldValue;
                this.setSelectionRange(this.oldSelectionStart, this.oldSelectionEnd);
            }
        });
    };

    var maxCeiling = $("#calc-val-ceiling").val();
    $("#calc-val").inputFilter(function (value) {
        return /^\d*$/.test(value) && (value === "" || parseInt(value) <= parseInt(maxCeiling))
    });

    $("#calc-button").click(function (event) {
        event.preventDefault();
        $(document.body).css({ "cursor": "wait" });
        $("#calc-result").html("Calculating...");
        var maxNum = $("#calc-val").val();
        var restList = $("#prime-link");
        $.get("/prime/" + maxNum, function (data) {
            // console.log(data);
            $(document.body).css({ "cursor": "default" });
            if (data && data.prime && data.prime.val) {
                $("#calc-result").html("Highest prime: <b>" + data.prime.val + "</b> (duration: " + data.dur + ")");
                restList.attr("href", "/prime/" + data.prime.max);
                restList.html("/prime/" + data.prime.max);
            } else {
                $("#calc-result").html("<b>Error:</b> " + data);
            }
        });
    });
});