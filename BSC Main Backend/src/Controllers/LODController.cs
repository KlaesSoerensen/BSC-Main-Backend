using BSC_Main_Backend.dto.response;
using Microsoft.AspNetCore.Mvc;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("[controller]")]
public class LODController : ControllerBase
{
    private readonly ILogger<LODController> _logger;
    
    public LODController(ILogger<LODController> logger)
    {
        _logger = logger;
    }
    
    /// <summary>
    /// Returns the LOD of that id.
    /// </summary>
    /// <param name="lodId">The ID of the LOD</param>
    /// <returns></returns>
    [HttpGet("{lodId}", Name = "GetLOD")]
    public AssetLODResponseDTO GetLod([FromRoute] uint lodId)
    {
        return null;
    }
    
    /// <summary>
    /// Get a specific LOD of some GraphicAsset of some detail level.
    /// </summary>
    /// <param name="assetId">Id of asset</param>
    /// <param name="detailLevel">Level of detail requested</param>
    /// <returns>That given LOD or 404 if no such LOD</returns>
    [HttpGet(Name = "GetLODByAsset")]
    public AssetLODResponseDTO GetLODByAsset([FromQuery] uint assetId, [FromQuery] uint detailLevel)
    {
        return null;
    }
    
    
}