using BSC_Main_Backend.dto.response;
using BSC_Main_Backend.dto.request;
using Microsoft.AspNetCore.Mvc;

namespace BSC_Main_Backend.Controllers;

[ApiController]
[Route("[controller]")]
public class PlayerController : ControllerBase
{
    
    private readonly ILogger<PlayerController> _logger;
    
    public PlayerController(ILogger<PlayerController> logger)
    {
        _logger = logger;
    }

    /// <summary>
    /// Get info about some player
    /// </summary>
    /// <param name="playerId">Id of said player</param>
    /// <returns></returns>
    [HttpGet("{playerId}", Name = "Get player info")]
    public PlayerInfoResponseDTO GetPlayerInfo([FromRoute] uint playerId)
    {
        return null;
    }

    /// <summary>
    /// Get player preferences
    /// </summary>
    /// <param name="playerId">Id of player</param>
    /// <returns></returns>
    [HttpGet("{playerId}/preferences", Name = "Get player preferences")]
    public PlayerPreferencesResponseDTO GetPlayerPreferences([FromRoute] uint playerId)
    {
        return null;
    }
    
    /// <summary>
    /// Get all colonies of some player
    /// </summary>
    /// <param name="playerId"></param>
    /// <returns></returns>
    [HttpGet("{playerId}/colonies", Name = "Get overview of all player colonies")]
    public ColonyOverviewResponseDTO GetAllColoniesOfPlayer([FromRoute] uint playerId)
    {
        return null;
    }
    
    /// <summary>
    /// Get a specific colony of some player
    /// </summary>
    /// <param name="playerId">Id of player</param>
    /// <param name="colonyId">Id of colony</param>
    /// <returns>That given colony info or 404 if no such colony/player</returns>
    [HttpGet("{playerId}/colony/{colonyId}", Name = "GetColonyOfPlayer")]
    public ColonyInfoResponseDTO GetColonyOfPlayer([FromRoute] uint playerId, [FromRoute] uint colonyId)
    {
        return null;
    }

    /// <summary>
    /// Create a new colony based on some settings.
    /// </summary>
    /// <param name="playerId">Id of player</param>
    /// <returns>The info pertaining to the newly created colony</returns>
    [HttpPost("{playerId}/colony/create", Name = "Create a new colony")]
    public ColonyInfoResponseDTO CreateANewColony([FromRoute] uint playerId,
        [FromBody] ColonyCreationRequestDTO creationDTO)
    {
        return null;
    }
    

}