$(document).ready(function() {
    $('[data-toggle="tooltip"]').tooltip()
    refreshDropdown();
    getPages()
});

$('#dropType').on('change', function() {
    $('#dropTypeTg').html($(this).val())
    getPages()
})

function refreshDropdown() {
    $('.select2-dropdown').select2();
}

function getPages() {
    var type = $('#dropType').val()
    if (type == 'GET') { $('#get-card').removeClass('d-none') } else { $('#get-card').addClass('d-none') }
    if (type == 'INSERT') { $('#insert-card').removeClass('d-none') } else { $('#insert-card').addClass('d-none') }
    if (type == 'UPDATE') { $('#update-card').removeClass('d-none') } else { $('#update-card').addClass('d-none') }
    if (type == 'DELETE') { $('#delete-card').removeClass('d-none') } else { $('#delete-card').addClass('d-none') }
}

$('#btn-simpan-api-modal').on('click', function(e) {
    e.preventDefault();
    $("#modal").modal('show');
});

$('#btn-simpan-api').on('click', function() {
    var endpoint = $('#input-endpoint').val()
    var koneksi = $('#input-koneksi').val()
    var query = $('#query-box').val()
    var type = $('#dropType').val()
    var data = $("input[name='jml_data']").val() * 1;
    var param = $("input[name='jml_param']").val() * 1;

    if (endpoint == '') {
        $('#input-endpoint').focus()
    } else if (koneksi == '') {
        $('#input-koneksi').focus()
    } else if (query == '') {
        $('#query-box').focus()
    } else {
        $.ajax({
            type: 'POST',
            url: "/save-endpoint",
            dataType: "json",
            data: {
                "Inpendpoint": endpoint,
                "Inptype": type,
                "Inpdata": data,
                "Inpparam": param,
                "Inpkoneksi": koneksi,
                "Inpquery": query,
            },
            async: false,
            success: function (data) {
                resetInput()
                $("pre").html('');
                $(".slide-form>.col-8").toggle();
                $(".box-tool>a,.box-tool>hr").remove();
                $("#modal").modal('hide');
                if (data.STATUS == "SUCCESS") {
                    $('.alert').removeClass("alert-success,alert-danger");
                    $('.alert').text("Data berhasil disimpan");
                    $('.alert').addClass("alert-success");
                    $('.alert').removeClass('d-none');
                    setTimeout(function() {
                        $('.alert').addClass('d-none');
                    }, 5000);
                    $("html, body").animate({scrollTop: 0}, 200);
                } else {
                    $('.alert').addClass("alert-danger");
                    $('.alert').text("Data tidak berhasil disimpan");
                    $('.alert').removeClass('d-none')
                    setTimeout(function() {
                        $('.alert').addClass('d-none')
                    }, 5000);
                    $("html, body").animate({scrollTop: 0}, 200);
                }
            },
            error: function(xhr, desc, err) {
                console.log(xhr);
                console.log("Details: " + desc + "\nError:" + err);
                $('#api-alert-failed').removeClass('d-none')
                setTimeout(function() {
                    $('#api-alert-failed').addClass('d-none')
                }, 5000);
            }
        });
    }
})

$('#btn-preview-api').on('click', function() {
    previewApiJson()
})

$('#btn-modal-api').on('click', function() {
    $('#modal-simpan-endpoint').modal('show');
    previewApiJson()   
})

function previewApiJson() {
    var endpoint = $('#input-endpoint').val()
    var koneksi = $('#input-koneksi').val()
    var query = $('#query-box').val()
    var type = $('#dropType').val()
    var data = 0;
    var param = 0;

    if (type == "GET") {
        param = ($('#get_params').val() != "") ? $('#get_params').val() : 0
    } else if (type == "INSERT") {
        data = ($('#insert_data').val() != "") ? $('#insert_data').val() : 0
    } else if (type == "UPDATE") {
        data = ($('#update_data').val() != "") ? $('#update_data').val() : 0
        param = ($('#update_params').val() != "") ? $('#update_params').val() : 0
    } else if (type == "DELETE") {
        param = ($('#delete_params').val() != "") ? $('#delete_params').val() : 0
    }

    if (endpoint == '') {
        $('#input-endpoint').focus()
    } else if (koneksi == '') {
        $('#input-koneksi').focus()
    } else if (query == '') {
        $('#query-box').focus()
    } else {
        $.ajax({
            type: 'POST',
            url: "/preview-endpoint",
            dataType: "json",
            data: {
                "Inpendpoint": endpoint,
                "Inptype": type,
                "Inpdata": data,
                "Inpparam": param,
                "Inpkoneksi": koneksi,
                "Inpquery": query,
            },
            async: false,
            success: function (data) {
                if (data.STATUS == "SUCCESS") {
                    $('#btn-simpan-api').show()
                    $('#json-data').text(JSON.stringify(data.JSON, undefined, 2))
                    $("html, body").animate({scrollTop: $(document).height()}, 200);
    
                    $('#json-data-modal').text(JSON.stringify(data.JSON, undefined, 2))
                    $("html, body").animate({scrollTop: $(document).height()}, 200);

                    $("#return-box").html("")
                    $("#return-box").html(data.HTML)
                } else {
                    $('#btn-simpan-api').hide()
                    $('#json-data').text(JSON.stringify(data, undefined, 2))
                    $("html, body").animate({scrollTop: $(document).height()}, 200);
    
                    $('#json-data-modal').text(JSON.stringify(data, undefined, 2))
                    $("html, body").animate({scrollTop: $(document).height()}, 200);
                    $("#return-box").html("")
                }
            },
            error: function(xhr, desc, err) {
                console.log(xhr);
                console.log("Details: " + desc + "\nError:" + err);
            }
        });
    }
}

function resetInput() {
    $('#input-endpoint').val('')
    $('#input-koneksi').val('ORCL1')
    $('#query-box').val('')

    $('#get_params').val('')
    $('#insert_data').val('')
    $('#update_data').val('')
    $('#update_params').val('')
    $('#delete_params').val('')
    
    $('#json-data').text('')
    $('#json-data-modal').text('')
}