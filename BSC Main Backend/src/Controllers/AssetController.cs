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

    [HttpGet("{assetId}",Name = "Get asset by id")]
    public AssetResponseDTO GetAsset([FromRoute] uint assetId)
    {
        return null;
    }
    
    [HttpGet("",Name = "Get multiple assets")]
    public AssetResponseDTO[] GetAssets([FromQuery] uint[] ids)
    {
        return null;
    }
}