using BSC_Main_Backend.dto;
using BSC_Main_Backend.dto.response;
using Microsoft.AspNetCore.Mvc;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("collection")]

public class AssetCollectionController : ControllerBase
{
    
    private readonly ILogger<AssetCollectionController> _logger;

    public AssetCollectionController(ILogger<AssetCollectionController> logger)
    {
        _logger = logger;
    }
    
    /// <summary>
    /// Returns the AssetCollection of that id.
    /// </summary>
    /// <param name="asssetCollectionId">The ID of the AssetCollection</param>
    /// <returns>AssetCollectionResponseDTO</returns>
    [HttpGet("{asssetCollectionId}",Name = "GetAssetCollection")]
    public AssetCollectionResponseDTO GetAssetCollection([FromRoute] uint collectionId)
    {
        // Fetch from database and perform logic as needed.
        var entries = new List<AssetCollectionEntryDTO>
        {
            // new AssetCollectionEntryDTO(1, new TransformDTO()),
            // new AssetCollectionEntryDTO(2, new TransformDTO())
        };

        var response = new AssetCollectionResponseDTO(
            Id: 123,
            Original: 456,
            Name: "Sample Collection",
            Entries: entries
        );

        return null; // response
    }
    
}