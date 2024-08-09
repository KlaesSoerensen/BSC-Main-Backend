using BSC_Main_Backend.dto;
using Microsoft.AspNetCore.Mvc;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("[controller]")]

public class AssetController
{
    private readonly ILogger<AssetController> _logger;

    public AssetController(ILogger<AssetController> logger)
    {
        _logger = logger;
    }

    [HttpGet(Name = "<base-URL>/asset/<assetId>")]
    public FetchAssetResponseDTO GetAsset()
    {
        return null;
    }
    
    [HttpGet(Name = "<base-URL>/assets?ids=x,y,z")]
    public FetchAssetResponseDTO GetAssets()
    {
        return null;
    }
}