using BSC_Main_Backend.dto.responses;
using BSC_Main_Backend.dto;
using Microsoft.AspNetCore.Mvc;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("[controller]")]

public class AssetCollectionController : ControllerBase
{
    
    private readonly ILogger<AssetCollectionController> _logger;

    public AssetCollectionController(ILogger<AssetCollectionController> logger)
    {
        _logger = logger;
    }
    
    [HttpGet(Name = "<base-URL>/collection/<collectionId>")]
    public FetchAssetCollectionResponseDTO GetAssetCollection()
    {
        // Fetch from database and perform logic as needed.
        var entries = new List<AssetCollectionEntryDTO>
        {
            // new AssetCollectionEntryDTO(1, new TransformDTO()),
            // new AssetCollectionEntryDTO(2, new TransformDTO())
        };

        var response = new FetchAssetCollectionResponseDTO(
            Id: 123,
            Original: 456,
            Name: "Sample Collection",
            Entries: entries
        );

        return null; // response
    }
    
}