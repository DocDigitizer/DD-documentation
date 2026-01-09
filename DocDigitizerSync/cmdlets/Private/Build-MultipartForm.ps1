function Build-MultipartForm {
    <#
    .SYNOPSIS
        Builds a multipart form data request for file upload.
    .DESCRIPTION
        Internal function that constructs the multipart/form-data content
        required for the DocIngester API.
    #>
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)]
        [string]$FilePath,

        [Parameter(Mandatory)]
        [guid]$DocumentId,

        [Parameter(Mandatory)]
        [guid]$ContextId,

        [Parameter()]
        [string]$RequestToken
    )

    $boundary = [System.Guid]::NewGuid().ToString()
    $fileName = [System.IO.Path]::GetFileName($FilePath)
    $fileBytes = [System.IO.File]::ReadAllBytes($FilePath)
    $fileEnc = [System.Text.Encoding]::GetEncoding('ISO-8859-1').GetString($fileBytes)

    $LF = "`r`n"

    $bodyLines = @(
        "--$boundary",
        "Content-Disposition: form-data; name=`"files`"; filename=`"$fileName`"",
        "Content-Type: application/pdf",
        "",
        $fileEnc,
        "--$boundary",
        "Content-Disposition: form-data; name=`"id`"",
        "",
        $DocumentId.ToString(),
        "--$boundary",
        "Content-Disposition: form-data; name=`"contextID`"",
        "",
        $ContextId.ToString()
    )

    if ($RequestToken) {
        $bodyLines += @(
            "--$boundary",
            "Content-Disposition: form-data; name=`"requestToken`"",
            "",
            $RequestToken
        )
    }

    $bodyLines += "--$boundary--"

    $body = $bodyLines -join $LF

    return @{
        Body        = $body
        ContentType = "multipart/form-data; boundary=$boundary"
    }
}
