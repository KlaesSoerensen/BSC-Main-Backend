using BSC_Main_Backend.dto;
using BSC_Main_Backend.dto.response;
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
    
    /// <summary>
    /// Returns the Asset of that id.
    /// </summary>
    /// <param name="assetId">The ID of the Asset</param>
    /// <returns>AssetResponseDTO</returns>
    [HttpGet("{assetId}", Name = "GetAsset")]
    public AssetResponseDTO GetAsset([FromRoute] uint assetId)
    {
        return null; // Implementation here
    }
    
    /// <summary>
    /// Returns the Assets of the id's.
    /// </summary>
    /// <param name="ids">The ID's of the Assets</param>
    /// <returns>AssetResponseDTO[]</returns>
    [HttpGet(Name = "GetAssets")]
    public AssetResponseDTO[] GetAssets([FromQuery] uint[] ids)
    {
        return null; // Implementation here
    }
}