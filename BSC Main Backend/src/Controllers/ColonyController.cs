using Microsoft.AspNetCore.Mvc;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("[controller]")]

public class ColonyController
{
    private readonly ILogger<ColonyController> _logger;

    public ColonyController(ILogger<ColonyController> logger)
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